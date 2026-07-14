package maven

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/danielvm-git/big-release/internal/publishers"
	"github.com/danielvm-git/big-release/internal/publishers/httputil"
)

const (
	DefaultRegistryURL = "https://central.sonatype.com/api/v1/publisher/upload"
	DefaultVerifyURL   = "https://search.maven.org/solrsearch/select"

	envToken        = "MAVEN_TOKEN"
	maxResponseSize = 10 * 1024 * 1024 // 10 MB
)

var validMavenIdentifier = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type pomProject struct {
	XMLName      xml.Name `xml:"project"`
	ModelVersion string   `xml:"modelVersion"`
	GroupID      string   `xml:"groupId"`
	ArtifactID   string   `xml:"artifactId"`
	Version      string   `xml:"version"`
	Packaging    string   `xml:"packaging"`
}

type Publisher struct {
	RegistryURL string
	Client      *httputil.RetryClient
	DryRun      bool
	VerifyURL   string
}

func NewPublisher() *Publisher {
	return &Publisher{
		RegistryURL: DefaultRegistryURL,
		Client:      httputil.NewRetryClient(http.DefaultClient),
		VerifyURL:   DefaultVerifyURL,
	}
}

func (p *Publisher) Name() string {
	return "maven"
}

func (p *Publisher) Detect() bool {
	_, err := os.Stat("pom.xml")
	return err == nil
}

func (p *Publisher) Prepare(version string) error {
	data, err := os.ReadFile("pom.xml")
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("maven: pom.xml not found")
		}
		return fmt.Errorf("maven: failed to read pom.xml: %w", err)
	}

	var pom pomProject
	if err := xml.Unmarshal(data, &pom); err != nil {
		return fmt.Errorf("maven: failed to parse pom.xml: %w", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	var out bytes.Buffer

	depth := 0
	inProjectVersion := false
	found := false

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("maven: failed to parse pom.xml: %w", err)
		}

		switch tok := t.(type) {
		case xml.StartElement:
			depth++
			if depth == 2 && tok.Name.Local == "version" {
				inProjectVersion = true
			}
			writeStartElement(&out, tok)
		case xml.EndElement:
			if depth == 2 && tok.Name.Local == "version" {
				inProjectVersion = false
			}
			depth--
			fmt.Fprintf(&out, "</%s>", tok.Name.Local)
		case xml.CharData:
			if inProjectVersion && !found {
				out.WriteString(version)
				found = true
			} else {
				out.Write(tok)
			}
		case xml.Comment:
			fmt.Fprintf(&out, "<!--%s-->", string(tok))
		case xml.ProcInst:
			fmt.Fprintf(&out, "<?%s %s?>", tok.Target, string(tok.Inst))
		case xml.Directive:
			fmt.Fprintf(&out, "<!%s>", string(tok))
		}
	}

	if !found {
		return fmt.Errorf("maven: version element not found in pom.xml")
	}

	if err := os.WriteFile("pom.xml", out.Bytes(), 0644); err != nil {
		return fmt.Errorf("maven: failed to write pom.xml: %w", err)
	}

	return nil
}

func writeStartElement(buf *bytes.Buffer, se xml.StartElement) {
	buf.WriteByte('<')
	buf.WriteString(se.Name.Local)

	type attrItem struct {
		name  string
		value string
	}
	var decls, attrs []attrItem

	for _, a := range se.Attr {
		if a.Name.Space == "xmlns" || (a.Name.Space == "" && a.Name.Local == "xmlns") {
			if a.Name.Space == "xmlns" {
				decls = append(decls, attrItem{name: "xmlns:" + a.Name.Local, value: a.Value})
			} else {
				decls = append(decls, attrItem{name: "xmlns", value: a.Value})
			}
		} else if a.Name.Space != "" && a.Name.Space != "xmlns" {
			attrs = append(attrs, attrItem{name: a.Name.Space + ":" + a.Name.Local, value: a.Value})
		} else {
			attrs = append(attrs, attrItem{name: a.Name.Local, value: a.Value})
		}
	}

	for _, d := range decls {
		buf.WriteByte(' ')
		buf.WriteString(d.name)
		buf.WriteString(`="`)
		writeEscaped(buf, d.value)
		buf.WriteByte('"')
	}
	for _, a := range attrs {
		buf.WriteByte(' ')
		buf.WriteString(a.name)
		buf.WriteString(`="`)
		writeEscaped(buf, a.value)
		buf.WriteByte('"')
	}

	buf.WriteByte('>')
}

func writeEscaped(buf *bytes.Buffer, s string) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			buf.WriteString("&quot;")
		case '&':
			buf.WriteString("&amp;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		default:
			buf.WriteByte(s[i])
		}
	}
}

func (p *Publisher) Publish(version string) error {
	token := os.Getenv(envToken)
	if token == "" {
		return fmt.Errorf("maven: %s environment variable is empty", envToken)
	}

	if p.DryRun {
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, p.RegistryURL, nil)
	if err != nil {
		return fmt.Errorf("maven: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("maven: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Drain body to allow connection reuse.
	_, _ = io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))

	return nil
}

func (p *Publisher) Verify(version string) error {
	groupID, artifactID, err := readPOMCoordinates()
	if err != nil {
		return fmt.Errorf("maven: %w", err)
	}

	if !isValidMavenIdentifier(groupID) {
		return fmt.Errorf("maven: invalid groupId %q", groupID)
	}
	if !isValidMavenIdentifier(artifactID) {
		return fmt.Errorf("maven: invalid artifactId %q", artifactID)
	}
	if !isValidMavenIdentifier(version) {
		return fmt.Errorf("maven: invalid version %q", version)
	}

	url := fmt.Sprintf("%s?q=g:%s+AND+a:%s+AND+v:%s&rows=1&wt=json",
		p.VerifyURL, groupID, artifactID, version)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("maven: failed to create verify request: %w", err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("maven: verify request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("maven: version %s not found (HTTP 404)", version)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("maven: verify failed with HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
	if err != nil {
		return fmt.Errorf("maven: failed to read verify response: %w", err)
	}

	var searchResponse struct {
		Response struct {
			NumFound int `json:"numFound"`
		} `json:"response"`
	}
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return fmt.Errorf("maven: failed to parse search response: %w", err)
	}
	if searchResponse.Response.NumFound == 0 {
		return fmt.Errorf("maven: version %s not found in Maven Central", version)
	}

	return nil
}

// SetDryRun sets the dry-run mode.
func (p *Publisher) SetDryRun(dryRun bool) {
	p.DryRun = dryRun
}

func init() {
	publishers.Register(NewPublisher())
}

func isValidMavenIdentifier(id string) bool {
	return len(id) > 0 && len(id) <= 256 && validMavenIdentifier.MatchString(id)
}

func readPOMCoordinates() (groupID, artifactID string, err error) {
	data, readErr := os.ReadFile("pom.xml")
	if readErr != nil {
		return "", "", fmt.Errorf("failed to read pom.xml: %w", readErr)
	}

	var pom pomProject
	if parseErr := xml.Unmarshal(data, &pom); parseErr != nil {
		return "", "", fmt.Errorf("failed to parse pom.xml: %w", parseErr)
	}

	if pom.GroupID == "" {
		return "", "", fmt.Errorf("groupId not found in pom.xml")
	}
	if pom.ArtifactID == "" {
		return "", "", fmt.Errorf("artifactId not found in pom.xml")
	}

	return pom.GroupID, pom.ArtifactID, nil
}
