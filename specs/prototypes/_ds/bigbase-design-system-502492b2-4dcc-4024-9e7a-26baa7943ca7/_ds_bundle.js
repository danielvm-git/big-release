/* @ds-bundle: {"format":4,"namespace":"BigBaseDesignSystem_502492","components":[{"name":"App","sourcePath":"source/bigbase-ui/App.tsx"},{"name":"Layout","sourcePath":"source/bigbase-ui/Layout.tsx"},{"name":"Avatar","sourcePath":"source/bigbase-ui/components/Avatar.tsx"},{"name":"Badge","sourcePath":"source/bigbase-ui/components/Badge.tsx"},{"name":"StatusBadge","sourcePath":"source/bigbase-ui/components/Badge.tsx"},{"name":"Button","sourcePath":"source/bigbase-ui/components/Button.tsx"},{"name":"Card","sourcePath":"source/bigbase-ui/components/Card.tsx"},{"name":"CardHeader","sourcePath":"source/bigbase-ui/components/Card.tsx"},{"name":"EmptyState","sourcePath":"source/bigbase-ui/components/EmptyState.tsx"},{"name":"Input","sourcePath":"source/bigbase-ui/components/Input.tsx"},{"name":"PageHeader","sourcePath":"source/bigbase-ui/components/PageHeader.tsx"},{"name":"Spinner","sourcePath":"source/bigbase-ui/components/Spinner.tsx"},{"name":"Tabs","sourcePath":"source/bigbase-ui/components/Tabs.tsx"},{"name":"ThemeSelect","sourcePath":"source/bigbase-ui/components/ThemeSelect.tsx"},{"name":"ThemeProvider","sourcePath":"source/bigbase-ui/theme/ThemeProvider.tsx"},{"name":"THEMES","sourcePath":"source/bigbase-ui/theme/themes.ts"},{"name":"THEME_BY_KEY","sourcePath":"source/bigbase-ui/theme/themes.ts"},{"name":"CiciPage","sourcePath":"source/bigbase-ui/pages/CiciPage.tsx"},{"name":"DashboardPage","sourcePath":"source/bigbase-ui/pages/DashboardPage.tsx"},{"name":"DataStudioPage","sourcePath":"source/bigbase-ui/pages/DataStudioPage.tsx"},{"name":"DeployPage","sourcePath":"source/bigbase-ui/pages/DeployPage.tsx"},{"name":"ForgePage","sourcePath":"source/bigbase-ui/pages/ForgePage.tsx"},{"name":"FunctionsPage","sourcePath":"source/bigbase-ui/pages/FunctionsPage.tsx"},{"name":"GitReposPage","sourcePath":"source/bigbase-ui/pages/GitReposPage.tsx"},{"name":"LoginPage","sourcePath":"source/bigbase-ui/pages/LoginPage.tsx"},{"name":"MessagingPage","sourcePath":"source/bigbase-ui/pages/MessagingPage.tsx"},{"name":"MonitoringPage","sourcePath":"source/bigbase-ui/pages/MonitoringPage.tsx"},{"name":"NotFoundPage","sourcePath":"source/bigbase-ui/pages/NotFoundPage.tsx"},{"name":"SqlEditorPage","sourcePath":"source/bigbase-ui/pages/SqlEditorPage.tsx"},{"name":"StoragePage","sourcePath":"source/bigbase-ui/pages/StoragePage.tsx"},{"name":"UsersPage","sourcePath":"source/bigbase-ui/pages/UsersPage.tsx"}],"sourceHashes":{"assets/icons.jsx":"ec0ddb29fb74","source/bigbase-ui/App.tsx":"ddce4adcd4e9","source/bigbase-ui/Layout.tsx":"36a443b760b6","source/bigbase-ui/components/Avatar.tsx":"b6f328072c22","source/bigbase-ui/components/Badge.tsx":"fd691d3c24d3","source/bigbase-ui/components/Button.tsx":"b9517a5cfebf","source/bigbase-ui/components/Card.tsx":"ed1d430f496c","source/bigbase-ui/components/EmptyState.tsx":"a9dcfe0546bd","source/bigbase-ui/components/Input.tsx":"b56b88345f72","source/bigbase-ui/components/PageHeader.tsx":"2d5f35b9f7dc","source/bigbase-ui/components/Spinner.tsx":"bee14b4707f3","source/bigbase-ui/components/Tabs.tsx":"b00d767925ca","source/bigbase-ui/components/ThemeSelect.tsx":"58f7d387c390","source/bigbase-ui/components/index.ts":"05d4094bff9d","source/bigbase-ui/main.tsx":"f4adb9ce2287","source/bigbase-ui/pages/CiciPage.tsx":"2c11aff9cd2f","source/bigbase-ui/pages/DashboardPage.tsx":"f193cf427e3c","source/bigbase-ui/pages/DataStudioPage.tsx":"9ef4351c47cc","source/bigbase-ui/pages/DeployPage.tsx":"3b30142be67a","source/bigbase-ui/pages/ForgePage.tsx":"6cb6ea2ebcd2","source/bigbase-ui/pages/FunctionsPage.tsx":"4cbd84aea5c9","source/bigbase-ui/pages/GitReposPage.tsx":"dd1b9ed622d7","source/bigbase-ui/pages/LoginPage.tsx":"9c16116d4aff","source/bigbase-ui/pages/MessagingPage.tsx":"463f30ea83ad","source/bigbase-ui/pages/MonitoringPage.tsx":"9fbfb58730c6","source/bigbase-ui/pages/NotFoundPage.tsx":"421af576ddd5","source/bigbase-ui/pages/SqlEditorPage.tsx":"46b5c4e93ca9","source/bigbase-ui/pages/StoragePage.tsx":"0ec2a50bd5d7","source/bigbase-ui/pages/UsersPage.tsx":"88431e177841","source/bigbase-ui/theme/ThemeProvider.tsx":"02de89fda121","source/bigbase-ui/theme/themes.ts":"d862cf9327d0","ui_kits/admin-console/app.jsx":"06fc2150b8b0","ui_kits/admin-console/data.jsx":"d07d7efb5b15","ui_kits/admin-console/screens.jsx":"61da27d112d8","ui_kits/admin-console/screens2.jsx":"0701f4e165d1","ui_kits/admin-console/ui.jsx":"96a19b1aa40e","ui_kits/admin-console/wizard.jsx":"a8ed6819761c"},"inlinedExternals":[],"unexposedExports":[{"name":"statusBadgeVariant","sourcePath":"source/bigbase-ui/components/Badge.tsx"},{"name":"useTheme","sourcePath":"source/bigbase-ui/theme/ThemeProvider.tsx"}]} */

(() => {

const __ds_ns = (window.BigBaseDesignSystem_502492 = window.BigBaseDesignSystem_502492 || {});

const __ds_scope = {};

(__ds_ns.__errors = __ds_ns.__errors || []);

// assets/icons.jsx
try { (() => {
/* BigBase Design System — Icon component
   Lucide-derived stroke icons (ISC licensed). Inlined so files
   stay self-contained. <Icon name="rocket" size={18} />
   Substitution note: Appwrite Console ships a proprietary icon
   font; BigBase adopts Lucide for an open, consistent 2px-stroke set. */
(function () {
  const P = {
    'layout-dashboard': '<rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/>',
    'database': '<ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5V19A9 3 0 0 0 21 19V5"/><path d="M3 12A9 3 0 0 0 21 12"/>',
    'terminal': '<polyline points="4 17 10 11 4 5"/><line x1="12" x2="20" y1="19" y2="19"/>',
    'folder': '<path d="M20 20a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.9a2 2 0 0 1-1.69-.9L9.6 3.9A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2Z"/>',
    'users': '<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>',
    'git-branch': '<line x1="6" x2="6" y1="3" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/>',
    'rocket': '<path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="m12 15-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.55-3.03 2-4c1.62-1.08 5 0 5 0"/><path d="M12 15v5s3.03-.55 4-2c1.08-1.62 0-5 0-5"/>',
    'code': '<polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/>',
    'send': '<path d="M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11z"/><path d="m21.854 2.147-10.94 10.939"/>',
    'box': '<path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/>',
    'git-pull-request': '<circle cx="18" cy="18" r="3"/><circle cx="6" cy="6" r="3"/><path d="M13 6h3a2 2 0 0 1 2 2v7"/><line x1="6" x2="6" y1="9" y2="21"/>',
    'activity': '<path d="M22 12h-2.48a2 2 0 0 0-1.93 1.46l-2.35 8.36a.25.25 0 0 1-.48 0L9.24 2.18a.25.25 0 0 0-.48 0l-2.35 8.36A2 2 0 0 1 4.49 12H2"/>',
    'plus': '<path d="M5 12h14"/><path d="M12 5v14"/>',
    'search': '<circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>',
    'check': '<path d="M20 6 9 17l-5-5"/>',
    'check-circle': '<path d="M21.801 10A10 10 0 1 1 17 3.335"/><path d="m9 11 3 3L22 4"/>',
    'chevron-right': '<path d="m9 18 6-6-6-6"/>',
    'chevron-left': '<path d="m15 18-6-6 6-6"/>',
    'chevron-down': '<path d="m6 9 6 6 6-6"/>',
    'external-link': '<path d="M15 3h6v6"/><path d="M10 14 21 3"/><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>',
    'refresh-cw': '<path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/><path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/><path d="M8 16H3v5"/>',
    'globe': '<circle cx="12" cy="12" r="10"/><path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20"/><path d="M2 12h20"/>',
    'github': '<path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 4 5 4 5 4c-.3 1.15-.3 2.35 0 3.5A5.4 5.4 0 0 0 4 11c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"/><path d="M9 18c-4.51 2-5-2-7-2"/>',
    'settings': '<path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/>',
    'log-out': '<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" x2="9" y1="12" y2="12"/>',
    'moon': '<path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"/>',
    'sun': '<circle cx="12" cy="12" r="4"/><path d="M12 2v2"/><path d="M12 20v2"/><path d="m4.93 4.93 1.41 1.41"/><path d="m17.66 17.66 1.41 1.41"/><path d="M2 12h2"/><path d="M20 12h2"/><path d="m6.34 17.66-1.41 1.41"/><path d="m19.07 4.93-1.41 1.41"/>',
    'x': '<path d="M18 6 6 18"/><path d="m6 6 12 12"/>',
    'alert-triangle': '<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/>',
    'info': '<circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/>',
    'clock': '<circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>',
    'more-horizontal': '<circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/>',
    'arrow-left': '<path d="m12 19-7-7 7-7"/><path d="M19 12H5"/>',
    'arrow-right': '<path d="M5 12h14"/><path d="m12 5 7 7-7 7"/>',
    'zap': '<path d="M4 14a1 1 0 0 1-.78-1.63l9.9-10.2a.5.5 0 0 1 .86.46l-1.92 6.02A1 1 0 0 0 13 10h7a1 1 0 0 1 .78 1.63l-9.9 10.2a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14z"/>',
    'hard-drive': '<line x1="22" x2="2" y1="12" y2="12"/><path d="M5.45 5.11 2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/><line x1="6" x2="6.01" y1="16" y2="16"/><line x1="10" x2="10.01" y1="16" y2="16"/>',
    'mail': '<rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>',
    'file-code': '<path d="M10 12.5 8 15l2 2.5"/><path d="m14 12.5 2 2.5-2 2.5"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M4 6V4a2 2 0 0 1 2-2h9l5 5v13a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-1"/>',
    'play': '<polygon points="6 3 20 12 6 21 6 3"/>',
    'copy': '<rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>',
    'link': '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>',
    'loader': '<line x1="12" x2="12" y1="2" y2="6"/><line x1="12" x2="12" y1="18" y2="22"/><line x1="4.93" x2="7.76" y1="4.93" y2="7.76"/><line x1="16.24" x2="19.07" y1="16.24" y2="19.07"/><line x1="2" x2="6" y1="12" y2="12"/><line x1="18" x2="22" y1="12" y2="12"/><line x1="4.93" x2="7.76" y1="19.07" y2="16.24"/><line x1="16.24" x2="19.07" y1="7.76" y2="4.93"/>',
    'cpu': '<rect width="16" height="16" x="4" y="4" rx="2"/><rect width="6" height="6" x="9" y="9" rx="1"/><path d="M15 2v2"/><path d="M15 20v2"/><path d="M2 15h2"/><path d="M2 9h2"/><path d="M20 15h2"/><path d="M20 9h2"/><path d="M9 2v2"/><path d="M9 20v2"/>',
    'bell': '<path d="M10.268 21a2 2 0 0 0 3.464 0"/><path d="M3.262 15.326A1 1 0 0 0 4 17h16a1 1 0 0 0 .74-1.673C19.41 13.956 18 12.499 18 8A6 6 0 0 0 6 8c0 4.499-1.411 5.956-2.738 7.326"/>'
  };
  function Icon({
    name,
    size = 18,
    className = '',
    style = {},
    strokeWidth = 2
  }) {
    const d = P[name] || '';
    return React.createElement('svg', {
      className: 'icon ' + className,
      width: size,
      height: size,
      viewBox: '0 0 24 24',
      fill: 'none',
      stroke: 'currentColor',
      strokeWidth,
      strokeLinecap: 'round',
      strokeLinejoin: 'round',
      style,
      dangerouslySetInnerHTML: {
        __html: d
      }
    });
  }
  window.Icon = Icon;
  window.ICON_PATHS = P;
})();
})(); } catch (e) { __ds_ns.__errors.push({ path: "assets/icons.jsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/Layout.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function SidebarIcon({
  children
}) {
  return /*#__PURE__*/React.createElement("span", {
    className: "sidebar-nav-icon"
  }, children);
}
function Layout() {
  const nav = useNavigate();
  const [user, setUser] = useState(null);
  useEffect(() => {
    fetch('/api/auth/me').then(r => {
      if (!r.ok) throw new Error('unauthorized');
      return r.json();
    }).then(setUser).catch(() => nav('/login'));
  }, [nav]);
  const handleLogout = () => nav('/login');
  const initial = user?.email?.[0]?.toUpperCase() || '?';
  return /*#__PURE__*/React.createElement("div", {
    className: "layout"
  }, /*#__PURE__*/React.createElement("nav", {
    className: "sidebar"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sidebar-logo"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sidebar-logo-icon"
  }, "B"), /*#__PURE__*/React.createElement("span", null, "BigBase")), /*#__PURE__*/React.createElement("div", {
    className: "sidebar-section"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sidebar-section-title"
  }, "Overview"), /*#__PURE__*/React.createElement("ul", {
    className: "sidebar-nav"
  }, /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/",
    end: true
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "H"), /*#__PURE__*/React.createElement("span", null, "Dashboard"))))), /*#__PURE__*/React.createElement("div", {
    className: "sidebar-section"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sidebar-section-title"
  }, "Data"), /*#__PURE__*/React.createElement("ul", {
    className: "sidebar-nav"
  }, /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/data"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "D"), /*#__PURE__*/React.createElement("span", null, "Data Studio"))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/sql"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "S"), /*#__PURE__*/React.createElement("span", null, "SQL Editor"))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/storage"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "F"), /*#__PURE__*/React.createElement("span", null, "Storage"))))), /*#__PURE__*/React.createElement("div", {
    className: "sidebar-section"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sidebar-section-title"
  }, "Services"), /*#__PURE__*/React.createElement("ul", {
    className: "sidebar-nav"
  }, /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/users"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "U"), /*#__PURE__*/React.createElement("span", null, "Users"))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/repos"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "G"), /*#__PURE__*/React.createElement("span", null, "Git Repos"))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/deploy"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "R"), /*#__PURE__*/React.createElement("span", null, "Deploy"))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/functions"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "\u03BB"), /*#__PURE__*/React.createElement("span", null, "Functions"))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/messaging"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "M"), /*#__PURE__*/React.createElement("span", null, "Messaging"))))), /*#__PURE__*/React.createElement("div", {
    className: "sidebar-section"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sidebar-section-title"
  }, "DevOps"), /*#__PURE__*/React.createElement("ul", {
    className: "sidebar-nav"
  }, /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/forge"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "I"), /*#__PURE__*/React.createElement("span", null, "Forge"))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/cici"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "C"), /*#__PURE__*/React.createElement("span", null, "CI/CD"))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(NavLink, {
    to: "/monitoring"
  }, /*#__PURE__*/React.createElement(SidebarIcon, null, "V"), /*#__PURE__*/React.createElement("span", null, "Monitoring"))))), /*#__PURE__*/React.createElement("div", {
    className: "sidebar-spacer"
  }), user && /*#__PURE__*/React.createElement("div", {
    className: "sidebar-footer"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sidebar-user"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sidebar-avatar"
  }, initial), /*#__PURE__*/React.createElement("span", {
    className: "sidebar-email"
  }, user.email)), /*#__PURE__*/React.createElement("button", {
    className: "btn btn-secondary btn-sm",
    onClick: handleLogout,
    style: {
      width: '100%'
    }
  }, "Logout"))), /*#__PURE__*/React.createElement("main", {
    className: "content"
  }, /*#__PURE__*/React.createElement(Outlet, null)));
}
Object.assign(__ds_scope, { Layout, __ds_default_source_bigbase_ui_Layout_2g8xmy: Layout });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/Layout.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/Avatar.tsx
try { (() => {
/** Round monogram avatar. Maps to `.avatar`; fill is the active accent. */

function Avatar({
  email,
  size = 30
}) {
  const initial = (email || '?')[0].toUpperCase();
  return /*#__PURE__*/React.createElement("div", {
    className: "avatar",
    style: {
      width: size,
      height: size,
      fontSize: size * 0.42
    }
  }, initial);
}
Object.assign(__ds_scope, { Avatar });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/Avatar.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/Badge.tsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
const variantClass = {
  neutral: 'badge-neutral',
  success: 'badge-success',
  warning: 'badge-warning',
  error: 'badge-error',
  accent: 'badge-accent',
  info: 'badge-info'
};
function Badge({
  variant = 'neutral',
  dot = false,
  className = '',
  children,
  ...rest
}) {
  return /*#__PURE__*/React.createElement("span", _extends({
    className: `badge ${variantClass[variant]} ${className}`
  }, rest), dot && /*#__PURE__*/React.createElement("span", {
    className: "dot"
  }), children);
}
function statusBadgeVariant(status) {
  const s = status.toLowerCase();
  if (s === 'running' || s === 'active' || s === 'ok' || s === 'healthy' || s === 'ready') return 'success';
  if (s === 'failed' || s === 'error' || s === 'deleted') return 'error';
  if (s === 'building' || s === 'pending' || s === 'deploying') return 'warning';
  return 'neutral';
}

/**
 * Status pill that pairs colour with a dot (or spinner while building) and a
 * Title-Case label — status is never communicated by colour alone (a11y).
 */
function StatusBadge({
  status,
  children
}) {
  const variant = statusBadgeVariant(status);
  const label = children ?? status.charAt(0).toUpperCase() + status.slice(1);
  const building = status.toLowerCase() === 'building';
  return /*#__PURE__*/React.createElement(Badge, {
    variant: variant
  }, building ? /*#__PURE__*/React.createElement("span", {
    className: "spinner spinner-sm",
    style: {
      width: 10,
      height: 10,
      borderWidth: 1.5
    }
  }) : /*#__PURE__*/React.createElement("span", {
    className: "dot"
  }), label);
}
Object.assign(__ds_scope, { Badge, statusBadgeVariant, StatusBadge });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/Badge.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/Button.tsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
const variantClass = {
  primary: 'btn-primary',
  secondary: 'btn-secondary',
  danger: 'btn-danger',
  ghost: 'btn-ghost',
  link: 'btn-link'
};
function Button(props) {
  const {
    variant = 'primary',
    size = 'md',
    loading = false,
    children,
    className = ''
  } = props;
  const classes = ['btn', variantClass[variant], size === 'sm' ? 'btn-sm' : '', className].filter(Boolean).join(' ');
  const content = loading ? /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("span", {
    className: "spinner spinner-sm"
  }), children) : children;
  if (props.as === 'a') {
    const {
      as: _,
      loading: _l,
      ...anchorProps
    } = props;
    return /*#__PURE__*/React.createElement("a", _extends({
      className: classes
    }, anchorProps), content);
  }
  const {
    as: _,
    loading: _l,
    ...buttonProps
  } = props;
  return /*#__PURE__*/React.createElement("button", _extends({
    className: classes,
    disabled: loading || buttonProps.disabled
  }, buttonProps), content);
}
Object.assign(__ds_scope, { Button });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/Button.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/Card.tsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
function Card({
  children,
  className = '',
  ...rest
}) {
  return /*#__PURE__*/React.createElement("div", _extends({
    className: `card ${className}`
  }, rest), children);
}
function CardHeader({
  title,
  children
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "card-header"
  }, /*#__PURE__*/React.createElement("span", {
    className: "card-title"
  }, title), children);
}
Object.assign(__ds_scope, { Card, CardHeader });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/Card.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/EmptyState.tsx
try { (() => {
function EmptyState({
  icon = '—',
  title,
  description,
  children
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "empty-state"
  }, /*#__PURE__*/React.createElement("span", {
    className: "empty-state-icon"
  }, icon), /*#__PURE__*/React.createElement("span", {
    className: "empty-state-title"
  }, title), description && /*#__PURE__*/React.createElement("span", {
    className: "empty-state-text"
  }, description), children);
}
Object.assign(__ds_scope, { EmptyState });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/EmptyState.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/Input.tsx
try { (() => {
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
function Input(props) {
  const {
    label,
    error,
    hint,
    className = '',
    ...rest
  } = props;
  const inputClass = `input ${error ? 'input-error' : ''} ${className}`.trim();
  const id = props.id || props.name;
  const inputElement = props.as === 'textarea' ? /*#__PURE__*/React.createElement("textarea", _extends({
    id: id,
    className: inputClass
  }, rest)) : props.as === 'select' ? /*#__PURE__*/React.createElement("select", _extends({
    id: id,
    className: inputClass
  }, rest), rest.children) : /*#__PURE__*/React.createElement("input", _extends({
    id: id,
    className: inputClass
  }, rest));
  return /*#__PURE__*/React.createElement("div", {
    className: "input-group"
  }, label && /*#__PURE__*/React.createElement("label", {
    htmlFor: id,
    className: "input-label"
  }, label), inputElement, hint && !error && /*#__PURE__*/React.createElement("span", {
    className: "input-hint"
  }, hint), error && /*#__PURE__*/React.createElement("span", {
    className: "input-error-text"
  }, error));
}
Object.assign(__ds_scope, { Input });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/Input.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/PageHeader.tsx
try { (() => {
function PageHeader({
  title,
  children
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "page-header"
  }, /*#__PURE__*/React.createElement("h1", {
    className: "page-title"
  }, title), children);
}
Object.assign(__ds_scope, { PageHeader });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/PageHeader.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/Spinner.tsx
try { (() => {
/** Indeterminate loading spinner. Maps to `.spinner` / `.spinner-sm`. */

function Spinner({
  size = 'md',
  className = ''
}) {
  return /*#__PURE__*/React.createElement("span", {
    className: `spinner ${size === 'sm' ? 'spinner-sm' : ''} ${className}`,
    role: "status",
    "aria-label": "Loading"
  });
}
Object.assign(__ds_scope, { Spinner });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/Spinner.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/Tabs.tsx
try { (() => {
function Tabs({
  tabs,
  active,
  onChange
}) {
  return /*#__PURE__*/React.createElement("div", {
    className: "tabs"
  }, tabs.map(tab => /*#__PURE__*/React.createElement("button", {
    key: tab.id,
    className: `tab ${active === tab.id ? 'active' : ''}`,
    onClick: () => onChange(tab.id)
  }, tab.label)));
}
Object.assign(__ds_scope, { Tabs });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/Tabs.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/theme/themes.ts
try { (() => {
/**
 * BigBase month themes — typed accent presets.
 * Mirrors the `[data-accent="…"]` blocks in colors_and_type.css.
 * Each theme is WCAG 2.1 AA verified (see SPECS.md §3).
 */

const THEMES = [{
  key: 'default',
  month: 'Default',
  label: 'Indigo',
  swatch: 'rgb(79, 70, 229)',
  onAccent: 'white'
}, {
  key: 'january',
  month: 'January',
  label: 'Teal',
  swatch: 'rgb(13, 148, 136)',
  onAccent: 'dark'
}, {
  key: 'february',
  month: 'February',
  label: 'Orange',
  swatch: 'rgb(234, 88, 12)',
  onAccent: 'dark'
}, {
  key: 'march',
  month: 'March',
  label: 'Purple',
  swatch: 'rgb(124, 58, 237)',
  onAccent: 'white'
}, {
  key: 'april',
  month: 'April',
  label: 'Green',
  swatch: 'rgb(22, 163, 74)',
  onAccent: 'dark'
}, {
  key: 'may',
  month: 'May',
  label: 'Lavender',
  swatch: 'rgb(167, 139, 250)',
  onAccent: 'dark'
}, {
  key: 'june',
  month: 'June',
  label: 'Rainbow',
  swatch: 'linear-gradient(to right, rgb(239,68,68), rgb(245,158,11), rgb(16,185,129), rgb(59,130,246), rgb(139,92,246))',
  onAccent: 'white'
}, {
  key: 'july',
  month: 'July',
  label: 'Peach',
  swatch: 'rgb(253, 186, 116)',
  onAccent: 'dark'
}, {
  key: 'august',
  month: 'August',
  label: 'Silver',
  swatch: 'rgb(156, 163, 175)',
  onAccent: 'dark'
}, {
  key: 'september',
  month: 'September',
  label: 'Yellow',
  swatch: 'rgb(234, 179, 8)',
  onAccent: 'dark'
}, {
  key: 'october',
  month: 'October',
  label: 'Pink',
  swatch: 'rgb(236, 72, 153)',
  onAccent: 'dark'
}, {
  key: 'november',
  month: 'November',
  label: 'Blue',
  swatch: 'rgb(37, 99, 235)',
  onAccent: 'white'
}, {
  key: 'december',
  month: 'December',
  label: 'Red',
  swatch: 'rgb(220, 38, 38)',
  onAccent: 'white'
}];
const THEME_BY_KEY = Object.fromEntries(THEMES.map(t => [t.key, t]));
Object.assign(__ds_scope, { THEMES, THEME_BY_KEY });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/theme/themes.ts", error: String((e && e.message) || e) }); }

// source/bigbase-ui/theme/ThemeProvider.tsx
try { (() => {
const {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback
} = React;
const THEME_STORAGE_KEY = 'bigbase-theme';
const DARK_STORAGE_KEY = 'bigbase-dark';
const ThemeContext = createContext(null);
function read(key, fallback) {
  try {
    return localStorage.getItem(key) ?? fallback;
  } catch {
    return fallback;
  }
}
/**
 * Applies `data-accent` and `data-theme` to <html> and persists to
 * localStorage. Components never read this directly — they consume the
 * CSS role tokens that these attributes remap. `useTheme()` is only for
 * UI that *changes* the theme (the sidebar ThemeSelect, the dark toggle).
 */
function ThemeProvider({
  children,
  defaultTheme,
  defaultDark
}) {
  const [theme, setThemeState] = useState(() => defaultTheme ?? read(THEME_STORAGE_KEY, 'default'));
  const [dark, setDarkState] = useState(() => defaultDark ?? read(DARK_STORAGE_KEY, 'false') === 'true');
  useEffect(() => {
    const root = document.documentElement;
    if (theme && theme !== 'default') root.setAttribute('data-accent', theme);else root.removeAttribute('data-accent');
    try {
      localStorage.setItem(THEME_STORAGE_KEY, theme);
    } catch {/* ignore */}
  }, [theme]);
  useEffect(() => {
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
    try {
      localStorage.setItem(DARK_STORAGE_KEY, String(dark));
    } catch {/* ignore */}
  }, [dark]);
  const setTheme = useCallback(key => setThemeState(key), []);
  const setDark = useCallback(on => setDarkState(on), []);
  const toggleDark = useCallback(() => setDarkState(d => !d), []);
  return /*#__PURE__*/React.createElement(ThemeContext.Provider, {
    value: {
      theme,
      setTheme,
      dark,
      setDark,
      toggleDark
    }
  }, children);
}
function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error('useTheme must be used within <ThemeProvider>');
  return ctx;
}
Object.assign(__ds_scope, { ThemeProvider, useTheme });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/theme/ThemeProvider.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/ThemeSelect.tsx
try { (() => {
const {
  useState,
  useRef,
  useEffect
} = React;
/**
 * Accent-theme dropdown for the sidebar footer. Reads/writes the
 * ThemeProvider; closes on Escape and outside-click. 12 month themes +
 * default indigo. Each option is WCAG-AA verified (SPECS.md §3).
 *
 * States: closed (default), open (menu, --ease-emphasized), item hover,
 * item selected (brand-tint), focus-visible (--focus-ring).
 */
function ThemeSelect({
  className = ''
}) {
  const {
    theme,
    setTheme
  } = __ds_scope.useTheme();
  const [open, setOpen] = useState(false);
  const ref = useRef(null);
  const current = __ds_scope.THEMES.find(t => t.key === theme) ?? __ds_scope.THEMES[0];
  useEffect(() => {
    if (!open) return;
    const onDown = e => {
      if (ref.current && !ref.current.contains(e.target)) setOpen(false);
    };
    const onEsc = e => {
      if (e.key === 'Escape') setOpen(false);
    };
    document.addEventListener('mousedown', onDown);
    document.addEventListener('keydown', onEsc);
    return () => {
      document.removeEventListener('mousedown', onDown);
      document.removeEventListener('keydown', onEsc);
    };
  }, [open]);
  const labelFor = (m, l) => m === 'Default' ? 'Indigo' : `${m} · ${l}`;
  return /*#__PURE__*/React.createElement("div", {
    className: `theme-select-wrap ${className}`,
    ref: ref
  }, /*#__PURE__*/React.createElement("button", {
    type: "button",
    className: "theme-select",
    "aria-haspopup": "listbox",
    "aria-expanded": open,
    onClick: () => setOpen(o => !o)
  }, /*#__PURE__*/React.createElement("span", {
    className: "theme-swatch",
    style: {
      background: current.swatch
    }
  }), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }, labelFor(current.month, current.label)), /*#__PURE__*/React.createElement("span", {
    "aria-hidden": true
  }, "\u25BE")), open && /*#__PURE__*/React.createElement("div", {
    className: "theme-menu",
    role: "listbox",
    "aria-label": "Accent theme"
  }, __ds_scope.THEMES.map(t => /*#__PURE__*/React.createElement("button", {
    key: t.key,
    type: "button",
    role: "option",
    "aria-selected": t.key === theme,
    className: `theme-menu-item ${t.key === theme ? 'selected' : ''}`,
    onClick: () => {
      setTheme(t.key);
      setOpen(false);
    }
  }, /*#__PURE__*/React.createElement("span", {
    className: "theme-dot",
    style: {
      background: t.swatch
    }
  }), /*#__PURE__*/React.createElement("span", {
    style: {
      flex: 1
    }
  }, t.month === 'Default' ? 'Indigo' : t.month), t.key === theme && /*#__PURE__*/React.createElement("span", {
    "aria-hidden": true
  }, "\u2713")))));
}
Object.assign(__ds_scope, { ThemeSelect });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/ThemeSelect.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/components/index.ts
try { (() => {

Object.assign(__ds_scope, { Button: __ds_scope.Button, Card: __ds_scope.Card, CardHeader: __ds_scope.CardHeader, Input: __ds_scope.Input, Badge: __ds_scope.Badge, StatusBadge: __ds_scope.StatusBadge, statusBadgeVariant: __ds_scope.statusBadgeVariant, Tabs: __ds_scope.Tabs, PageHeader: __ds_scope.PageHeader, EmptyState: __ds_scope.EmptyState, Spinner: __ds_scope.Spinner, Avatar: __ds_scope.Avatar, ThemeSelect: __ds_scope.ThemeSelect, ThemeProvider: __ds_scope.ThemeProvider, useTheme: __ds_scope.useTheme, THEMES: __ds_scope.THEMES, THEME_BY_KEY: __ds_scope.THEME_BY_KEY });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/components/index.ts", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/CiciPage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function CiciPage() {
  const [repos, setRepos] = useState([]);
  const [repoId, setRepoId] = useState('');
  const [workflows, setWorkflows] = useState([]);
  const [runs, setRuns] = useState([]);
  const [logs, setLogs] = useState({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({
    name: '',
    yaml: ''
  });
  const [expandedRun, setExpandedRun] = useState(null);
  useEffect(() => {
    fetch('/api/git/repos').then(r => r.json()).then(d => setRepos(d.data || [])).catch(() => {});
  }, []);
  useEffect(() => {
    if (!repoId) return;
    setLoading(true);
    setError('');
    Promise.all([fetch(`/api/cici/${repoId}/workflows`).then(r => r.json()), fetch(`/api/cici/runs?repo_id=${repoId}`).then(r => r.json())]).then(([wD, rD]) => {
      setWorkflows(wD.data || []);
      setRuns(rD.data || []);
    }).catch(() => setError('failed to load')).finally(() => setLoading(false));
  }, [repoId]);
  const loadLogs = async runId => {
    if (expandedRun === runId) {
      setExpandedRun(null);
      return;
    }
    setExpandedRun(runId);
    try {
      const res = await fetch(`/api/cici/runs/${runId}/logs`);
      if (res.ok) {
        const d = await res.json();
        setLogs(p => ({
          ...p,
          [runId]: d.logs || []
        }));
      }
    } catch {}
  };
  const handleCreateWorkflow = async e => {
    e.preventDefault();
    setError('');
    try {
      const res = await fetch(`/api/cici/${repoId}/workflows`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(form)
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'save failed');
        return;
      }
      setShowForm(false);
      setForm({
        name: '',
        yaml: ''
      });
      const wRes = await fetch(`/api/cici/${repoId}/workflows`);
      const wD = await wRes.json();
      setWorkflows(wD.data || []);
    } catch {
      setError('network error');
    }
  };
  const handleRun = async wfId => {
    setError('');
    try {
      const res = await fetch(`/api/cici/${repoId}/workflows/${wfId}/run`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          event: 'manual'
        })
      });
      if (!res.ok) {
        const d = await res.json();
        setError(d.error || 'run failed');
        return;
      }
      const rRes = await fetch(`/api/cici/runs?repo_id=${repoId}`);
      const rD = await rRes.json();
      setRuns(rD.data || []);
    } catch {
      setError('network error');
    }
  };
  const ciciTabs = [{
    id: 'workflows',
    label: 'Workflows'
  }, {
    id: 'runs',
    label: 'Runs'
  }];
  const [ciciTab, setCiciTab] = useState('workflows');
  if (!repoId) return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "CI/CD"
  }), /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "Select a repo to view workflows."), /*#__PURE__*/React.createElement("select", {
    value: repoId,
    onChange: e => setRepoId(e.target.value),
    className: "input",
    style: {
      maxWidth: 240,
      marginTop: 'var(--space-4)'
    }
  }, /*#__PURE__*/React.createElement("option", {
    value: ""
  }, "Select repo..."), repos.map(r => /*#__PURE__*/React.createElement("option", {
    key: r.id,
    value: r.id
  }, r.name))));
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading...");
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "CI/CD"
  }, /*#__PURE__*/React.createElement("select", {
    value: repoId,
    onChange: e => setRepoId(e.target.value),
    className: "input",
    style: {
      maxWidth: 200
    }
  }, repos.map(r => /*#__PURE__*/React.createElement("option", {
    key: r.id,
    value: r.id
  }, r.name))), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "primary",
    size: "sm",
    onClick: () => setShowForm(!showForm)
  }, showForm ? 'Cancel' : 'New Workflow')), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), showForm && /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("h3", {
    style: {
      marginBottom: 'var(--space-6)',
      fontSize: 'var(--text-m)',
      fontWeight: 600
    }
  }, "New Workflow"), /*#__PURE__*/React.createElement("form", {
    onSubmit: handleCreateWorkflow,
    className: "fn-form"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Workflow name *",
    value: form.name,
    onChange: e => setForm(p => ({
      ...p,
      name: e.target.value
    })),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "textarea",
    placeholder: "YAML config *",
    value: form.yaml,
    onChange: e => setForm(p => ({
      ...p,
      yaml: e.target.value
    })),
    required: true,
    rows: 10,
    className: "code-textarea"
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit"
  }, "Save"))), /*#__PURE__*/React.createElement(__ds_scope.Tabs, {
    tabs: ciciTabs,
    active: ciciTab,
    onChange: setCiciTab
  }), ciciTab === 'workflows' && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("h2", {
    className: "section-title"
  }, "Workflows"), workflows.length === 0 ? /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No workflows.") : /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Actions"))), /*#__PURE__*/React.createElement("tbody", null, workflows.map(w => /*#__PURE__*/React.createElement("tr", {
    key: w.id
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, w.name)), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: () => handleRun(w.id)
  }, "Run")))))))), ciciTab === 'runs' && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("h2", {
    className: "section-title"
  }, "Runs"), runs.length === 0 ? /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No runs yet.") : /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "ID"), /*#__PURE__*/React.createElement("th", null, "Event"), /*#__PURE__*/React.createElement("th", null, "Status"), /*#__PURE__*/React.createElement("th", null, "Started"), /*#__PURE__*/React.createElement("th", null, "Finished"), /*#__PURE__*/React.createElement("th", null, "Logs"))), /*#__PURE__*/React.createElement("tbody", null, runs.map(run => /*#__PURE__*/React.createElement("tr", {
    key: run.id
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, run.id.slice(0, 8))), /*#__PURE__*/React.createElement("td", null, run.event), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Badge, {
    variant: __ds_scope.statusBadgeVariant(run.status)
  }, run.status)), /*#__PURE__*/React.createElement("td", null, run.started_at ? new Date(run.started_at).toLocaleString() : '—'), /*#__PURE__*/React.createElement("td", null, run.finished_at ? new Date(run.finished_at).toLocaleString() : '—'), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: () => loadLogs(run.id)
  }, expandedRun === run.id ? 'Hide' : 'Logs')))))))), expandedRun && logs[expandedRun] && /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginTop: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("h3", {
    style: {
      marginBottom: 'var(--space-4)'
    }
  }, "Logs \u2014 ", /*#__PURE__*/React.createElement("code", null, expandedRun.slice(0, 8))), logs[expandedRun].length === 0 ? /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No logs.") : /*#__PURE__*/React.createElement("div", {
    className: "code-output"
  }, logs[expandedRun].map((l, i) => /*#__PURE__*/React.createElement("div", {
    key: i,
    className: "log-entry"
  }, /*#__PURE__*/React.createElement("span", {
    className: "log-step"
  }, "[", l.step, "]"), /*#__PURE__*/React.createElement("pre", null, l.output))))));
}
Object.assign(__ds_scope, { CiciPage, __ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c: CiciPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/CiciPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/DashboardPage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function fmtSize(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
function DashboardPage() {
  const [user, setUser] = useState(null);
  const [health, setHealth] = useState(null);
  const [deployments, setDeployments] = useState([]);
  const [messages, setMessages] = useState([]);
  const [files, setFiles] = useState([]);
  const [stats, setStats] = useState([]);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const ctrl = new AbortController();
    const opts = {
      signal: ctrl.signal
    };
    Promise.all([fetch('/api/auth/me', opts).then(r => r.ok ? r.json() : null).catch(() => null), fetch('/health', opts).then(r => r.json()).catch(() => ({
      status: 'unknown',
      components: 0
    })), fetch('/api/git/repos', opts).then(r => r.json()).catch(() => ({
      data: []
    })), fetch('/api/deploy', opts).then(r => r.json()).catch(() => ({
      data: []
    })), fetch('/api/messaging/messages', opts).then(r => r.json()).catch(() => ({
      data: []
    })), fetch('/api/storage/files', opts).then(r => r.json()).catch(() => ({
      data: []
    })), fetch('/api/functions', opts).then(r => r.json()).catch(() => ({
      data: []
    }))]).then(([u, h, reposR, depR, msgR, fileR, fnR]) => {
      setUser(u);
      setHealth(h);
      const deps = depR.data || [];
      const msgs = msgR.data || [];
      const fils = fileR.data || [];
      setDeployments(deps);
      setMessages(msgs);
      setFiles(fils);
      setStats([{
        label: 'Git Repos',
        count: (reposR.data || []).length,
        link: '/repos'
      }, {
        label: 'Deployments',
        count: deps.length,
        link: '/deploy'
      }, {
        label: 'Messages',
        count: msgs.length,
        link: '/messaging'
      }, {
        label: 'Files',
        count: fils.length,
        link: '/storage'
      }, {
        label: 'Functions',
        count: (fnR.data || []).length,
        link: '/functions'
      }]);
    }).catch(() => {}).finally(() => setLoading(false));
    return () => ctrl.abort();
  }, []);
  const depByStatus = {};
  deployments.forEach(d => {
    depByStatus[d.status] = (depByStatus[d.status] || 0) + 1;
  });
  const depTotal = deployments.length;
  const msgByChannel = {};
  messages.forEach(m => {
    msgByChannel[m.channel] = (msgByChannel[m.channel] || 0) + 1;
  });
  const msgTotal = messages.length;
  const totalBytes = files.reduce((acc, f) => acc + f.size, 0);
  const statusColors = {
    running: '#22c55e',
    failed: '#ef4444',
    building: '#f59e0b',
    pending: '#6b7280'
  };
  const channelColors = {
    email: '#3b82f6',
    sms: '#8b5cf6',
    push: '#ec4899'
  };
  const recentDeps = deployments.slice(0, 5);
  const recentMsgs = messages.slice(0, 5);
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading dashboard...");
  if (!user) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading...");
  return /*#__PURE__*/React.createElement("div", {
    className: "dashboard"
  }, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Dashboard"
  }), /*#__PURE__*/React.createElement("div", {
    className: "dash-grid"
  }, /*#__PURE__*/React.createElement(__ds_scope.Card, null, /*#__PURE__*/React.createElement(__ds_scope.CardHeader, {
    title: "Signed in as"
  }), /*#__PURE__*/React.createElement("p", {
    style: {
      fontSize: 'var(--text-m)',
      fontWeight: 500
    }
  }, user.email), /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "ID #", user.id)), health && /*#__PURE__*/React.createElement(__ds_scope.Card, null, /*#__PURE__*/React.createElement(__ds_scope.CardHeader, {
    title: "System"
  }), /*#__PURE__*/React.createElement("p", {
    className: "stat-value"
  }, health.components), /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "components \xB7 ", health.status)), /*#__PURE__*/React.createElement(__ds_scope.Card, null, /*#__PURE__*/React.createElement(__ds_scope.CardHeader, {
    title: "Storage"
  }), /*#__PURE__*/React.createElement("p", {
    className: "stat-value"
  }, fmtSize(totalBytes)), /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, files.length, " files"))), stats.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "stats-grid"
  }, stats.map(s => /*#__PURE__*/React.createElement("a", {
    key: s.label,
    href: `/admin/#${s.link}`,
    className: "stat-card"
  }, /*#__PURE__*/React.createElement("span", {
    className: "stat-count"
  }, s.count), /*#__PURE__*/React.createElement("span", {
    className: "stat-label"
  }, s.label)))), /*#__PURE__*/React.createElement("div", {
    className: "dash-cols"
  }, /*#__PURE__*/React.createElement("div", {
    className: "dash-col"
  }, depTotal > 0 && /*#__PURE__*/React.createElement(__ds_scope.Card, null, /*#__PURE__*/React.createElement(__ds_scope.CardHeader, {
    title: "Deployments by Status"
  }), /*#__PURE__*/React.createElement("div", {
    className: "bar-chart"
  }, Object.entries(depByStatus).map(([status, count]) => /*#__PURE__*/React.createElement("div", {
    key: status,
    className: "bar-row"
  }, /*#__PURE__*/React.createElement("span", {
    className: "bar-label"
  }, status), /*#__PURE__*/React.createElement("div", {
    className: "bar-track"
  }, /*#__PURE__*/React.createElement("div", {
    className: "bar-fill",
    style: {
      width: `${count / depTotal * 100}%`,
      background: statusColors[status] || '#6b7280'
    }
  })), /*#__PURE__*/React.createElement("span", {
    className: "bar-count"
  }, count))))), msgTotal > 0 && /*#__PURE__*/React.createElement(__ds_scope.Card, null, /*#__PURE__*/React.createElement(__ds_scope.CardHeader, {
    title: "Messages by Channel"
  }), /*#__PURE__*/React.createElement("div", {
    className: "bar-chart"
  }, Object.entries(msgByChannel).map(([ch, count]) => /*#__PURE__*/React.createElement("div", {
    key: ch,
    className: "bar-row"
  }, /*#__PURE__*/React.createElement("span", {
    className: "bar-label"
  }, ch), /*#__PURE__*/React.createElement("div", {
    className: "bar-track"
  }, /*#__PURE__*/React.createElement("div", {
    className: "bar-fill",
    style: {
      width: `${count / msgTotal * 100}%`,
      background: channelColors[ch] || '#6b7280'
    }
  })), /*#__PURE__*/React.createElement("span", {
    className: "bar-count"
  }, count)))))), /*#__PURE__*/React.createElement("div", {
    className: "dash-col"
  }, recentDeps.length > 0 && /*#__PURE__*/React.createElement(__ds_scope.Card, null, /*#__PURE__*/React.createElement(__ds_scope.CardHeader, {
    title: "Recent Deployments"
  }), /*#__PURE__*/React.createElement("div", {
    className: "activity-feed"
  }, recentDeps.map(d => /*#__PURE__*/React.createElement("div", {
    key: d.id,
    className: "activity-item"
  }, /*#__PURE__*/React.createElement(__ds_scope.Badge, {
    variant: __ds_scope.statusBadgeVariant(d.status)
  }, d.status), /*#__PURE__*/React.createElement("span", {
    className: "activity-text"
  }, /*#__PURE__*/React.createElement("code", null, d.repo_id.slice(0, 8)), " \xB7 ", d.app_type || '?'), /*#__PURE__*/React.createElement("span", {
    className: "activity-time"
  }, new Date(d.created_at).toLocaleDateString()))))), recentMsgs.length > 0 && /*#__PURE__*/React.createElement(__ds_scope.Card, null, /*#__PURE__*/React.createElement(__ds_scope.CardHeader, {
    title: "Recent Messages"
  }), /*#__PURE__*/React.createElement("div", {
    className: "activity-feed"
  }, recentMsgs.map(m => /*#__PURE__*/React.createElement("div", {
    key: m.id,
    className: "activity-item"
  }, /*#__PURE__*/React.createElement("span", {
    className: `channel-badge channel-${m.channel}`
  }, m.channel), /*#__PURE__*/React.createElement("span", {
    className: "activity-text"
  }, m.to_addr), /*#__PURE__*/React.createElement("span", {
    className: "activity-time"
  }, new Date(m.created_at).toLocaleDateString()))))))));
}
Object.assign(__ds_scope, { DashboardPage, __ds_default_source_bigbase_ui_pages_DashboardPage_120zej4: DashboardPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/DashboardPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/DataStudioPage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function DataStudioPage() {
  const [collections, setCollections] = useState([]);
  const [selected, setSelected] = useState(null);
  const [records, setRecords] = useState([]);
  const [loading, setLoading] = useState(true);
  const [recordError, setRecordError] = useState('');
  useEffect(() => {
    fetch('/api/collections/').then(r => r.json()).then(d => setCollections(d.data || [])).catch(() => setCollections([])).finally(() => setLoading(false));
  }, []);
  const loadRecords = async name => {
    setSelected(name);
    setRecordError('');
    try {
      const res = await fetch(`/api/collections/${name}`);
      if (!res.ok) {
        setRecordError(`error: ${res.status}`);
        setRecords([]);
        return;
      }
      const d = await res.json();
      setRecords(d.data || []);
    } catch {
      setRecordError('network error');
      setRecords([]);
    }
  };
  const allKeys = records.length > 0 ? Array.from(new Set(records.flatMap(r => Object.keys(r)))) : [];
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading collections...");
  return /*#__PURE__*/React.createElement("div", {
    className: "data-studio"
  }, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Data Studio"
  }), /*#__PURE__*/React.createElement("div", {
    className: "studio-layout"
  }, /*#__PURE__*/React.createElement("aside", {
    className: "collection-list"
  }, /*#__PURE__*/React.createElement("div", {
    className: "collection-list-title"
  }, "Collections"), collections.length === 0 && /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No collections yet."), /*#__PURE__*/React.createElement("ul", {
    className: "collection-list-nav"
  }, collections.map(c => /*#__PURE__*/React.createElement("li", {
    key: c
  }, /*#__PURE__*/React.createElement("button", {
    className: `collection-btn${selected === c ? ' active' : ''}`,
    onClick: () => loadRecords(c)
  }, c))))), /*#__PURE__*/React.createElement("section", {
    className: "record-view"
  }, !selected && /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "Select a collection to browse."), recordError && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, recordError), selected && !recordError && records.length === 0 && /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No records found."), selected && records.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, allKeys.map(k => /*#__PURE__*/React.createElement("th", {
    key: k
  }, k)))), /*#__PURE__*/React.createElement("tbody", null, records.map(r => /*#__PURE__*/React.createElement("tr", {
    key: r.id
  }, allKeys.map(k => /*#__PURE__*/React.createElement("td", {
    className: "max",
    key: k
  }, JSON.stringify(r[k] ?? '')))))))))));
}
Object.assign(__ds_scope, { DataStudioPage, __ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6: DataStudioPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/DataStudioPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/DeployPage.tsx
try { (() => {
const {
  useEffect,
  useState,
  useCallback
} = React;
function DeployPage() {
  const [deployments, setDeployments] = useState([]);
  const [repos, setRepos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [selectedRepoId, setSelectedRepoId] = useState('');
  const [branch, setBranch] = useState('main');
  const fetchDeployments = useCallback(async () => {
    try {
      const res = await fetch('/api/deploy');
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || `error: ${res.status}`);
      } else {
        setDeployments(d.data || []);
      }
    } catch {
      setError('network error');
    } finally {
      setLoading(false);
    }
  }, []);
  const fetchRepos = useCallback(async () => {
    try {
      const res = await fetch('/api/git/repos');
      const d = await res.json();
      if (res.ok) {
        setRepos(d.data || []);
      }
    } catch {}
  }, []);
  useEffect(() => {
    fetchDeployments();
    fetchRepos();
  }, [fetchDeployments, fetchRepos]);
  useEffect(() => {
    const hasActive = deployments.some(d => d.status === 'pending' || d.status === 'building');
    if (!hasActive) return;
    const timer = setInterval(fetchDeployments, 3000);
    return () => clearInterval(timer);
  }, [deployments, fetchDeployments]);
  const handleCreate = async e => {
    e.preventDefault();
    if (!selectedRepoId) return;
    setError('');
    try {
      const res = await fetch('/api/deploy', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          repo_id: selectedRepoId,
          branch
        })
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'create failed');
        return;
      }
      setShowForm(false);
      setSelectedRepoId('');
      setBranch('main');
      fetchDeployments();
    } catch {
      setError('network error');
    }
  };
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading deployments...");
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Deployments"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: fetchDeployments
  }, "Refresh"), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "primary",
    size: "sm",
    onClick: () => setShowForm(!showForm)
  }, showForm ? 'Cancel' : 'New Deployment')), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), showForm && /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("form", {
    onSubmit: handleCreate,
    className: "form-row"
  }, /*#__PURE__*/React.createElement("select", {
    className: "input",
    value: selectedRepoId,
    onChange: e => setSelectedRepoId(e.target.value),
    required: true,
    style: {
      flex: 1,
      minWidth: 140
    }
  }, /*#__PURE__*/React.createElement("option", {
    value: ""
  }, "Select repo..."), repos.map(r => /*#__PURE__*/React.createElement("option", {
    key: r.id,
    value: r.id
  }, r.name))), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Branch",
    value: branch,
    onChange: e => setBranch(e.target.value)
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    size: "sm"
  }, "Deploy"))), deployments.length === 0 && !error && /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No deployments yet."), deployments.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Status"), /*#__PURE__*/React.createElement("th", null, "Repo"), /*#__PURE__*/React.createElement("th", null, "Branch"), /*#__PURE__*/React.createElement("th", null, "Type"), /*#__PURE__*/React.createElement("th", null, "URL"), /*#__PURE__*/React.createElement("th", null, "Commit"), /*#__PURE__*/React.createElement("th", null, "Created"))), /*#__PURE__*/React.createElement("tbody", null, deployments.map(d => /*#__PURE__*/React.createElement("tr", {
    key: d.id
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Badge, {
    variant: __ds_scope.statusBadgeVariant(d.status)
  }, d.status)), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, d.repo_id.slice(0, 8))), /*#__PURE__*/React.createElement("td", null, d.branch), /*#__PURE__*/React.createElement("td", null, d.app_type || '—'), /*#__PURE__*/React.createElement("td", null, d.url ? /*#__PURE__*/React.createElement("a", {
    href: d.url,
    target: "_blank",
    rel: "noreferrer"
  }, d.url) : '—'), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, d.commit_sha ? d.commit_sha.slice(0, 7) : '—')), /*#__PURE__*/React.createElement("td", null, new Date(d.created_at).toLocaleString())))))));
}
Object.assign(__ds_scope, { DeployPage, __ds_default_source_bigbase_ui_pages_DeployPage_17gbuph: DeployPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/DeployPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/ForgePage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function ForgePage() {
  const [repos, setRepos] = useState([]);
  const [repoId, setRepoId] = useState('');
  const [tab, setTab] = useState('issues');
  const [issues, setIssues] = useState([]);
  const [labels, setLabels] = useState([]);
  const [board, setBoard] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({
    title: '',
    description: '',
    status: 'open',
    labels: ''
  });
  const [selectedIssue, setSelectedIssue] = useState(null);
  const [comments, setComments] = useState([]);
  const [commentText, setCommentText] = useState('');
  useEffect(() => {
    fetch('/api/git/repos').then(r => r.json()).then(d => setRepos(d.data || [])).catch(() => {});
  }, []);
  useEffect(() => {
    if (!repoId) return;
    setLoading(true);
    setError('');
    Promise.all([fetch(`/api/forge/issues?repo_id=${repoId}`).then(r => r.json()), fetch(`/api/forge/labels?repo_id=${repoId}`).then(r => r.json()), fetch(`/api/forge/board?repo_id=${repoId}`).then(r => r.json())]).then(([iD, lD, bD]) => {
      setIssues(iD.data || []);
      setLabels(lD.data || []);
      setBoard(bD);
    }).catch(() => setError('failed to load')).finally(() => setLoading(false));
  }, [repoId]);
  const loadComments = async issueId => {
    try {
      const res = await fetch(`/api/forge/issues/${issueId}/comments`);
      if (res.ok) {
        const d = await res.json();
        setComments(d);
      }
    } catch {}
  };
  const handleCreate = async e => {
    e.preventDefault();
    setError('');
    try {
      const res = await fetch('/api/forge/issues', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          ...form,
          repo_id: repoId
        })
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'create failed');
        return;
      }
      setShowForm(false);
      setForm({
        title: '',
        description: '',
        status: 'open',
        labels: ''
      });
      const iRes = await fetch(`/api/forge/issues?repo_id=${repoId}`);
      const iD = await iRes.json();
      setIssues(iD.data || []);
    } catch {
      setError('network error');
    }
  };
  const handleStatusChange = async (id, status) => {
    try {
      await fetch(`/api/forge/issues/${id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          status
        })
      });
      const iRes = await fetch(`/api/forge/issues?repo_id=${repoId}`);
      const iD = await iRes.json();
      setIssues(iD.data || []);
      if (board) {
        const bRes = await fetch(`/api/forge/board?repo_id=${repoId}`);
        setBoard(await bRes.json());
      }
    } catch {}
  };
  const handleComment = async e => {
    e.preventDefault();
    if (!selectedIssue || !commentText.trim()) return;
    try {
      const res = await fetch(`/api/forge/issues/${selectedIssue.id}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          content: commentText
        })
      });
      if (res.ok) {
        setCommentText('');
        loadComments(selectedIssue.id);
      }
    } catch {}
  };
  if (!repoId) return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Forge"
  }), /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "Select a repo to view issues."), /*#__PURE__*/React.createElement("select", {
    value: repoId,
    onChange: e => setRepoId(e.target.value),
    className: "input",
    style: {
      maxWidth: 240,
      marginTop: 'var(--space-4)'
    }
  }, /*#__PURE__*/React.createElement("option", {
    value: ""
  }, "Select repo..."), repos.map(r => /*#__PURE__*/React.createElement("option", {
    key: r.id,
    value: r.id
  }, r.name))));
  const labelColor = name => labels.find(l => l.name === name)?.color || '#888';
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading...");
  const forgeTabs = [{
    id: 'issues',
    label: 'Issues'
  }, {
    id: 'board',
    label: 'Board'
  }];
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Forge"
  }, /*#__PURE__*/React.createElement("select", {
    value: repoId,
    onChange: e => setRepoId(e.target.value),
    className: "input",
    style: {
      maxWidth: 200
    }
  }, repos.map(r => /*#__PURE__*/React.createElement("option", {
    key: r.id,
    value: r.id
  }, r.name))), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "primary",
    size: "sm",
    onClick: () => setShowForm(!showForm)
  }, showForm ? 'Cancel' : 'New Issue')), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), /*#__PURE__*/React.createElement(__ds_scope.Tabs, {
    tabs: forgeTabs,
    active: tab,
    onChange: setTab
  }), showForm && /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("h3", {
    style: {
      marginBottom: 'var(--space-6)',
      fontSize: 'var(--text-m)',
      fontWeight: 600
    }
  }, "New Issue"), /*#__PURE__*/React.createElement("form", {
    onSubmit: handleCreate,
    className: "fn-form"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Title *",
    value: form.title,
    onChange: e => setForm(p => ({
      ...p,
      title: e.target.value
    })),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "textarea",
    placeholder: "Description",
    value: form.description,
    onChange: e => setForm(p => ({
      ...p,
      description: e.target.value
    })),
    rows: 4
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Labels (comma-separated)",
    value: form.labels,
    onChange: e => setForm(p => ({
      ...p,
      labels: e.target.value
    }))
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit"
  }, "Create"))), tab === 'issues' && /*#__PURE__*/React.createElement(React.Fragment, null, issues.length === 0 ? /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No issues.") : /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Title"), /*#__PURE__*/React.createElement("th", null, "Status"), /*#__PURE__*/React.createElement("th", null, "Labels"), /*#__PURE__*/React.createElement("th", null, "Updated"), /*#__PURE__*/React.createElement("th", null, "Actions"))), /*#__PURE__*/React.createElement("tbody", null, issues.map(issue => /*#__PURE__*/React.createElement("tr", {
    key: issue.id
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("button", {
    className: "btn-link",
    onClick: () => {
      setSelectedIssue(issue);
      loadComments(issue.id);
    }
  }, issue.title)), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("select", {
    value: issue.status,
    onChange: e => handleStatusChange(issue.id, e.target.value),
    className: "status-select"
  }, /*#__PURE__*/React.createElement("option", {
    value: "open"
  }, "Open"), /*#__PURE__*/React.createElement("option", {
    value: "in_progress"
  }, "In Progress"), /*#__PURE__*/React.createElement("option", {
    value: "review"
  }, "Review"), /*#__PURE__*/React.createElement("option", {
    value: "closed"
  }, "Closed"))), /*#__PURE__*/React.createElement("td", null, issue.labels ? issue.labels.split(',').map(l => /*#__PURE__*/React.createElement("span", {
    key: l,
    className: "label-badge",
    style: {
      background: labelColor(l.trim())
    }
  }, l.trim())) : '—'), /*#__PURE__*/React.createElement("td", null, new Date(issue.updated_at).toLocaleString()), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: () => handleStatusChange(issue.id, issue.status === 'closed' ? 'open' : 'closed')
  }, issue.status === 'closed' ? 'Reopen' : 'Close'))))))), selectedIssue && /*#__PURE__*/React.createElement("div", {
    className: "card issue-detail"
  }, /*#__PURE__*/React.createElement("div", {
    className: "issue-detail-header"
  }, /*#__PURE__*/React.createElement("h3", null, selectedIssue.title), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: () => setSelectedIssue(null)
  }, "Close")), /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, selectedIssue.description || 'No description.'), comments.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "comments"
  }, /*#__PURE__*/React.createElement("h4", {
    style: {
      marginBottom: 'var(--space-4)'
    }
  }, "Comments"), comments.map(c => /*#__PURE__*/React.createElement("div", {
    key: c.id,
    className: "comment"
  }, /*#__PURE__*/React.createElement("p", null, c.content), /*#__PURE__*/React.createElement("span", {
    className: "dim"
  }, new Date(c.created_at).toLocaleString())))), /*#__PURE__*/React.createElement("form", {
    onSubmit: handleComment,
    className: "form-row",
    style: {
      marginTop: 'var(--space-4)'
    }
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Add a comment...",
    value: commentText,
    onChange: e => setCommentText(e.target.value)
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    size: "sm"
  }, "Comment")))), tab === 'board' && board && /*#__PURE__*/React.createElement("div", {
    className: "board"
  }, ['open', 'in_progress', 'review', 'closed'].map(col => /*#__PURE__*/React.createElement("div", {
    key: col,
    className: "board-col"
  }, /*#__PURE__*/React.createElement("h3", {
    className: "board-col-title"
  }, col.replace('_', ' ')), (board[col] || []).map(issue => /*#__PURE__*/React.createElement("div", {
    key: issue.id,
    className: "board-card",
    onClick: () => {
      setSelectedIssue(issue);
      setTab('issues');
      loadComments(issue.id);
    }
  }, /*#__PURE__*/React.createElement("strong", null, issue.title), issue.labels && /*#__PURE__*/React.createElement("div", {
    className: "board-labels"
  }, issue.labels.split(',').map(l => /*#__PURE__*/React.createElement("span", {
    key: l,
    className: "label-badge",
    style: {
      background: labelColor(l.trim())
    }
  }, l.trim())))))))));
}
Object.assign(__ds_scope, { ForgePage, __ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf: ForgePage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/ForgePage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/FunctionsPage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
const defaultFn = {
  name: '',
  runtime: 'javascript',
  source: '',
  trigger: 'http',
  schedule: '',
  env: '{}',
  timeout: 30
};
function FunctionsPage() {
  const [fns, setFns] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form, setForm] = useState(defaultFn);
  const [runResult, setRunResult] = useState({});
  const fetchFns = async () => {
    setLoading(true);
    try {
      const res = await fetch('/api/functions');
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || `error: ${res.status}`);
        setFns([]);
      } else {
        setFns(d.data || []);
      }
    } catch {
      setError('network error');
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    fetchFns();
  }, []);
  const openCreate = () => {
    setEditing(null);
    setForm(defaultFn);
    setShowForm(true);
  };
  const openEdit = fn => {
    setEditing(fn);
    setForm({
      name: fn.name,
      runtime: fn.runtime,
      source: fn.source,
      trigger: fn.trigger,
      schedule: fn.schedule,
      env: fn.env,
      timeout: fn.timeout
    });
    setShowForm(true);
  };
  const handleSave = async e => {
    e.preventDefault();
    setError('');
    const isNew = !editing;
    const url = isNew ? '/api/functions' : `/api/functions/${editing.id}`;
    const method = isNew ? 'POST' : 'PUT';
    try {
      const res = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(form)
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'save failed');
        return;
      }
      setShowForm(false);
      fetchFns();
    } catch {
      setError('network error');
    }
  };
  const handleDelete = async id => {
    if (!confirm('Delete this function?')) return;
    try {
      const res = await fetch(`/api/functions/${id}`, {
        method: 'DELETE'
      });
      if (!res.ok) {
        const d = await res.json();
        setError(d.error || 'delete failed');
        return;
      }
      setFns(prev => prev.filter(f => f.id !== id));
    } catch {
      setError('network error');
    }
  };
  const handleRun = async id => {
    setError('');
    try {
      const res = await fetch(`/api/functions/${id}/run`, {
        method: 'POST'
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'run failed');
        return;
      }
      setRunResult(p => ({
        ...p,
        [id]: d
      }));
    } catch {
      setError('network error');
    }
  };
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading functions...");
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Functions"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: fetchFns
  }, "Refresh"), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "primary",
    size: "sm",
    onClick: openCreate
  }, "New Function")), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), showForm && /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("h3", {
    style: {
      marginBottom: 'var(--space-6)',
      fontSize: 'var(--text-m)',
      fontWeight: 600
    }
  }, editing ? 'Edit Function' : 'New Function'), /*#__PURE__*/React.createElement("form", {
    onSubmit: handleSave,
    className: "fn-form"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Name *",
    value: form.name,
    onChange: e => setForm(p => ({
      ...p,
      name: e.target.value
    })),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "select",
    value: form.runtime,
    onChange: e => setForm(p => ({
      ...p,
      runtime: e.target.value
    }))
  }, /*#__PURE__*/React.createElement("option", {
    value: "javascript"
  }, "JavaScript")), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "select",
    value: form.trigger,
    onChange: e => setForm(p => ({
      ...p,
      trigger: e.target.value
    }))
  }, /*#__PURE__*/React.createElement("option", {
    value: "http"
  }, "HTTP"), /*#__PURE__*/React.createElement("option", {
    value: "schedule"
  }, "Schedule"), /*#__PURE__*/React.createElement("option", {
    value: "event"
  }, "Event")), form.trigger === 'schedule' && /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Cron schedule",
    value: form.schedule,
    onChange: e => setForm(p => ({
      ...p,
      schedule: e.target.value
    }))
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Env JSON ({\"KEY\":\"val\"})",
    value: form.env,
    onChange: e => setForm(p => ({
      ...p,
      env: e.target.value
    }))
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Timeout (seconds)",
    type: "number",
    value: form.timeout,
    onChange: e => setForm(p => ({
      ...p,
      timeout: +e.target.value
    }))
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "textarea",
    placeholder: "Source code *",
    value: form.source,
    onChange: e => setForm(p => ({
      ...p,
      source: e.target.value
    })),
    required: true,
    rows: 8,
    className: "code-textarea"
  }), /*#__PURE__*/React.createElement("div", {
    className: "form-actions"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit"
  }, editing ? 'Update' : 'Create'), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "button",
    variant: "secondary",
    onClick: () => setShowForm(false)
  }, "Cancel")))), fns.length === 0 && !error && /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No functions yet."), fns.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Runtime"), /*#__PURE__*/React.createElement("th", null, "Trigger"), /*#__PURE__*/React.createElement("th", null, "Timeout"), /*#__PURE__*/React.createElement("th", null, "Created"), /*#__PURE__*/React.createElement("th", null, "Actions"))), /*#__PURE__*/React.createElement("tbody", null, fns.map(fn => /*#__PURE__*/React.createElement("tr", {
    key: fn.id
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, fn.name)), /*#__PURE__*/React.createElement("td", null, fn.runtime), /*#__PURE__*/React.createElement("td", null, fn.trigger), /*#__PURE__*/React.createElement("td", null, fn.timeout, "s"), /*#__PURE__*/React.createElement("td", null, new Date(fn.created_at).toLocaleString()), /*#__PURE__*/React.createElement("td", {
    className: "actions-cell"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: () => handleRun(fn.id)
  }, "Run"), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: () => openEdit(fn)
  }, "Edit"), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "danger",
    size: "sm",
    onClick: () => handleDelete(fn.id)
  }, "Delete"))))))), Object.entries(runResult).map(([id, result]) => /*#__PURE__*/React.createElement("div", {
    key: id,
    className: "card",
    style: {
      marginTop: 'var(--space-8)',
      maxHeight: 400,
      overflow: 'auto'
    }
  }, /*#__PURE__*/React.createElement("h3", {
    style: {
      marginBottom: 'var(--space-4)'
    }
  }, "Run Result \u2014 ", /*#__PURE__*/React.createElement("code", null, fns.find(f => f.id === id)?.name || id)), result.logs?.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "code-output"
  }, result.logs.map((l, i) => /*#__PURE__*/React.createElement("pre", {
    key: i
  }, l))), result.error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, result.error), result.result !== undefined && /*#__PURE__*/React.createElement("div", {
    className: "code-output",
    style: {
      marginTop: 'var(--space-4)'
    }
  }, /*#__PURE__*/React.createElement("pre", null, JSON.stringify(result.result, null, 2))))));
}
Object.assign(__ds_scope, { FunctionsPage, __ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx: FunctionsPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/FunctionsPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/GitReposPage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function GitReposPage() {
  const [repos, setRepos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({
    name: '',
    description: '',
    private: false
  });
  const fetchRepos = async () => {
    setError('');
    setLoading(true);
    try {
      const res = await fetch('/api/git/repos');
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || `error: ${res.status}`);
        setRepos([]);
      } else {
        setRepos(d.data || []);
      }
    } catch {
      setError('network error');
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    fetchRepos();
  }, []);
  const handleCreate = async e => {
    e.preventDefault();
    try {
      const res = await fetch('/api/git/repos', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(form)
      });
      if (!res.ok) {
        const d = await res.json();
        setError(d.error || 'create failed');
        return;
      }
      setShowForm(false);
      setForm({
        name: '',
        description: '',
        private: false
      });
      fetchRepos();
    } catch {
      setError('network error');
    }
  };
  const handleDelete = async id => {
    if (!confirm('Delete this repo?')) return;
    try {
      const res = await fetch(`/api/git/repos/${id}`, {
        method: 'DELETE'
      });
      if (!res.ok) {
        const d = await res.json();
        setError(d.error || 'delete failed');
        return;
      }
      setRepos(prev => prev.filter(r => r.id !== id));
    } catch {
      setError('network error');
    }
  };
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading repos...");
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Git Repos"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: fetchRepos
  }, "Refresh"), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "primary",
    size: "sm",
    onClick: () => setShowForm(!showForm)
  }, showForm ? 'Cancel' : 'New Repo')), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), showForm && /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("form", {
    onSubmit: handleCreate,
    className: "form-row"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Name *",
    value: form.name,
    onChange: e => setForm(p => ({
      ...p,
      name: e.target.value
    })),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Description",
    value: form.description,
    onChange: e => setForm(p => ({
      ...p,
      description: e.target.value
    }))
  }), /*#__PURE__*/React.createElement("label", {
    className: "checkbox-label"
  }, /*#__PURE__*/React.createElement("input", {
    type: "checkbox",
    checked: form.private,
    onChange: e => setForm(p => ({
      ...p,
      private: e.target.checked
    }))
  }), "Private"), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    size: "sm"
  }, "Create"))), repos.length === 0 && !error && /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No repos yet."), repos.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Branch"), /*#__PURE__*/React.createElement("th", null, "Description"), /*#__PURE__*/React.createElement("th", null, "Created"), /*#__PURE__*/React.createElement("th", null, "Actions"))), /*#__PURE__*/React.createElement("tbody", null, repos.map(r => /*#__PURE__*/React.createElement("tr", {
    key: r.id
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, r.name)), /*#__PURE__*/React.createElement("td", null, r.default_branch), /*#__PURE__*/React.createElement("td", {
    className: "dim"
  }, r.description || '—'), /*#__PURE__*/React.createElement("td", null, new Date(r.created_at).toLocaleString()), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "danger",
    size: "sm",
    onClick: () => handleDelete(r.id)
  }, "Delete"))))))));
}
Object.assign(__ds_scope, { GitReposPage, __ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5: GitReposPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/GitReposPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/LoginPage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function LoginPage() {
  const nav = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isRegister, setIsRegister] = useState(false);
  const [googleEnabled, setGoogleEnabled] = useState(false);
  useEffect(() => {
    fetch('/api/auth/oauth/google', {
      redirect: 'manual'
    }).then(r => {
      if (r.status === 302) setGoogleEnabled(true);
    }).catch(() => {});
  }, []);
  const handleSubmit = async e => {
    e.preventDefault();
    setError('');
    const endpoint = isRegister ? '/api/auth/register' : '/api/auth/login';
    try {
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email,
          password
        })
      });
      const data = await res.json();
      if (!res.ok) {
        setError(data.error || 'request failed');
        return;
      }
      nav('/');
    } catch {
      setError('network error');
    }
  };
  return /*#__PURE__*/React.createElement("div", {
    className: "login-page"
  }, /*#__PURE__*/React.createElement(__ds_scope.Card, {
    className: "login-card"
  }, /*#__PURE__*/React.createElement("div", {
    className: "login-brand"
  }, /*#__PURE__*/React.createElement("div", {
    className: "login-brand-logo"
  }, "B"), /*#__PURE__*/React.createElement("h1", null, "BigBase"), /*#__PURE__*/React.createElement("p", null, isRegister ? 'Create your account' : 'Sign in to continue')), /*#__PURE__*/React.createElement("form", {
    onSubmit: handleSubmit,
    className: "login-form"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    type: "email",
    placeholder: "Email",
    value: email,
    onChange: e => setEmail(e.target.value),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    type: "password",
    placeholder: "Password",
    value: password,
    onChange: e => setPassword(e.target.value),
    required: true,
    minLength: 6
  }), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    variant: "primary"
  }, isRegister ? 'Register' : 'Sign In')), googleEnabled && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
    className: "divider"
  }, /*#__PURE__*/React.createElement("span", null, "or")), /*#__PURE__*/React.createElement("a", {
    href: "/api/auth/oauth/google",
    className: "google-btn"
  }, /*#__PURE__*/React.createElement("svg", {
    viewBox: "0 0 24 24",
    width: "18",
    height: "18"
  }, /*#__PURE__*/React.createElement("path", {
    fill: "#4285F4",
    d: "M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"
  }), /*#__PURE__*/React.createElement("path", {
    fill: "#34A853",
    d: "M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
  }), /*#__PURE__*/React.createElement("path", {
    fill: "#FBBC05",
    d: "M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
  }), /*#__PURE__*/React.createElement("path", {
    fill: "#EA4335",
    d: "M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
  })), "Sign in with Google")), /*#__PURE__*/React.createElement("p", {
    className: "login-toggle"
  }, isRegister ? 'Already have an account?' : "Don't have an account?", ' ', /*#__PURE__*/React.createElement("button", {
    className: "btn-link",
    onClick: () => setIsRegister(!isRegister),
    type: "button"
  }, isRegister ? 'Sign In' : 'Register'))));
}
Object.assign(__ds_scope, { LoginPage, __ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1: LoginPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/LoginPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/MessagingPage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function MessagingPage() {
  const [tab, setTab] = useState('email');
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [sending, setSending] = useState(false);
  const [emailTo, setEmailTo] = useState('');
  const [emailSubject, setEmailSubject] = useState('');
  const [emailBody, setEmailBody] = useState('');
  const [smsTo, setSmsTo] = useState('');
  const [smsMsg, setSmsMsg] = useState('');
  const [pushToken, setPushToken] = useState('');
  const [pushTitle, setPushTitle] = useState('');
  const [pushBody, setPushBody] = useState('');
  const fetchMessages = async () => {
    setLoading(true);
    try {
      const res = await fetch('/api/messaging/messages');
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || `error: ${res.status}`);
        setMessages([]);
      } else {
        setMessages(d.data || []);
      }
    } catch {
      setError('network error');
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    fetchMessages();
  }, []);
  const sendEmail = async e => {
    e.preventDefault();
    setSending(true);
    setError('');
    try {
      const res = await fetch('/api/messaging/email', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          to: emailTo,
          subject: emailSubject,
          body: emailBody
        })
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'send failed');
        setSending(false);
        return;
      }
      setEmailTo('');
      setEmailSubject('');
      setEmailBody('');
      fetchMessages();
    } catch {
      setError('network error');
    } finally {
      setSending(false);
    }
  };
  const sendSms = async e => {
    e.preventDefault();
    setSending(true);
    setError('');
    try {
      const res = await fetch('/api/messaging/sms', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          to: smsTo,
          message: smsMsg
        })
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'send failed');
        setSending(false);
        return;
      }
      setSmsTo('');
      setSmsMsg('');
      fetchMessages();
    } catch {
      setError('network error');
    } finally {
      setSending(false);
    }
  };
  const sendPush = async e => {
    e.preventDefault();
    setSending(true);
    setError('');
    try {
      const res = await fetch('/api/messaging/push', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          token: pushToken,
          title: pushTitle,
          body: pushBody
        })
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'send failed');
        setSending(false);
        return;
      }
      setPushToken('');
      setPushTitle('');
      setPushBody('');
      fetchMessages();
    } catch {
      setError('network error');
    } finally {
      setSending(false);
    }
  };
  const channelTabs = [{
    id: 'email',
    label: 'EMAIL'
  }, {
    id: 'sms',
    label: 'SMS'
  }, {
    id: 'push',
    label: 'PUSH'
  }];
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Messaging"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: fetchMessages
  }, "Refresh")), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), /*#__PURE__*/React.createElement(__ds_scope.Tabs, {
    tabs: channelTabs,
    active: tab,
    onChange: id => setTab(id)
  }), /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-12)'
    }
  }, tab === 'email' && /*#__PURE__*/React.createElement("form", {
    onSubmit: sendEmail,
    className: "msg-form"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "To",
    value: emailTo,
    onChange: e => setEmailTo(e.target.value),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Subject",
    value: emailSubject,
    onChange: e => setEmailSubject(e.target.value)
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "textarea",
    placeholder: "Body",
    value: emailBody,
    onChange: e => setEmailBody(e.target.value),
    required: true,
    rows: 4
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    disabled: sending
  }, sending ? 'Sending...' : 'Send Email')), tab === 'sms' && /*#__PURE__*/React.createElement("form", {
    onSubmit: sendSms,
    className: "msg-form"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "To",
    value: smsTo,
    onChange: e => setSmsTo(e.target.value),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "textarea",
    placeholder: "Message",
    value: smsMsg,
    onChange: e => setSmsMsg(e.target.value),
    required: true,
    rows: 3
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    disabled: sending
  }, sending ? 'Sending...' : 'Send SMS')), tab === 'push' && /*#__PURE__*/React.createElement("form", {
    onSubmit: sendPush,
    className: "msg-form"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Device Token",
    value: pushToken,
    onChange: e => setPushToken(e.target.value),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Title",
    value: pushTitle,
    onChange: e => setPushTitle(e.target.value)
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "textarea",
    placeholder: "Body",
    value: pushBody,
    onChange: e => setPushBody(e.target.value),
    required: true,
    rows: 3
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    disabled: sending
  }, sending ? 'Sending...' : 'Send Push'))), /*#__PURE__*/React.createElement("h2", {
    className: "section-title"
  }, "History"), loading ? /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading...") : messages.length === 0 ? /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No messages sent.") : /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Channel"), /*#__PURE__*/React.createElement("th", null, "To"), /*#__PURE__*/React.createElement("th", null, "Subject"), /*#__PURE__*/React.createElement("th", null, "Status"), /*#__PURE__*/React.createElement("th", null, "Sent"))), /*#__PURE__*/React.createElement("tbody", null, messages.map(m => /*#__PURE__*/React.createElement("tr", {
    key: m.id
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("span", {
    className: `channel-badge channel-${m.channel}`
  }, m.channel)), /*#__PURE__*/React.createElement("td", null, m.to_addr), /*#__PURE__*/React.createElement("td", {
    className: "dim"
  }, m.subject || '—'), /*#__PURE__*/React.createElement("td", null, m.status), /*#__PURE__*/React.createElement("td", null, new Date(m.created_at).toLocaleString())))))));
}
Object.assign(__ds_scope, { MessagingPage, __ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza: MessagingPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/MessagingPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/MonitoringPage.tsx
try { (() => {
const {
  useEffect,
  useState,
  useCallback
} = React;
function MonitoringPage() {
  const [metrics, setMetrics] = useState(null);
  const [logs, setLogs] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [tab, setTab] = useState('overview');
  const [logQuery, setLogQuery] = useState('');
  const [showAlertForm, setShowAlertForm] = useState(false);
  const [alertForm, setAlertForm] = useState({
    name: '',
    metric: '',
    threshold: 0,
    operator: 'gt',
    enabled: true
  });
  const fetchMetrics = useCallback(async () => {
    try {
      const res = await fetch('/api/monitoring/metrics');
      if (res.ok) setMetrics(await res.json());
    } catch {}
  }, []);
  const fetchLogs = useCallback(async q => {
    try {
      const url = q ? `/api/monitoring/logs?q=${encodeURIComponent(q)}` : '/api/monitoring/logs';
      const res = await fetch(url);
      if (res.ok) {
        const d = await res.json();
        setLogs(d.data || []);
      }
    } catch {}
  }, []);
  const fetchAlerts = useCallback(async () => {
    try {
      const res = await fetch('/api/monitoring/alerts');
      if (res.ok) {
        const d = await res.json();
        setAlerts(d.data || []);
      }
    } catch {}
  }, []);
  useEffect(() => {
    setLoading(true);
    Promise.all([fetchMetrics(), fetchLogs(), fetchAlerts()]).finally(() => setLoading(false));
  }, [fetchMetrics, fetchLogs, fetchAlerts]);
  useEffect(() => {
    if (!metrics) return;
    const timer = setInterval(fetchMetrics, 5000);
    return () => clearInterval(timer);
  }, [metrics, fetchMetrics]);
  const handleCreateAlert = async e => {
    e.preventDefault();
    try {
      const res = await fetch('/api/monitoring/alerts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(alertForm)
      });
      if (!res.ok) {
        const d = await res.json();
        setError(d.error || 'create failed');
        return;
      }
      setShowAlertForm(false);
      setAlertForm({
        name: '',
        metric: '',
        threshold: 0,
        operator: 'gt',
        enabled: true
      });
      fetchAlerts();
    } catch {
      setError('network error');
    }
  };
  const fmtUptime = s => {
    const d = Math.floor(s / 86400);
    const h = Math.floor(s % 86400 / 3600);
    return `${d}d ${h}h`;
  };
  const monitoringTabs = [{
    id: 'overview',
    label: 'Overview'
  }, {
    id: 'logs',
    label: 'Logs'
  }, {
    id: 'alerts',
    label: 'Alerts'
  }];
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading monitoring...");
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Monitoring"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: () => {
      fetchMetrics();
      fetchLogs();
      fetchAlerts();
    }
  }, "Refresh")), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), /*#__PURE__*/React.createElement(__ds_scope.Tabs, {
    tabs: monitoringTabs,
    active: tab,
    onChange: setTab
  }), tab === 'overview' && metrics && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("h2", {
    className: "section-title"
  }, "System"), /*#__PURE__*/React.createElement("div", {
    className: "stats-grid"
  }, /*#__PURE__*/React.createElement("div", {
    className: "stat-card"
  }, /*#__PURE__*/React.createElement("span", {
    className: "stat-count"
  }, metrics.system.memory_mb.toFixed(1)), /*#__PURE__*/React.createElement("span", {
    className: "stat-label"
  }, "Memory MB")), /*#__PURE__*/React.createElement("div", {
    className: "stat-card"
  }, /*#__PURE__*/React.createElement("span", {
    className: "stat-count"
  }, metrics.system.goroutines), /*#__PURE__*/React.createElement("span", {
    className: "stat-label"
  }, "Goroutines")), /*#__PURE__*/React.createElement("div", {
    className: "stat-card"
  }, /*#__PURE__*/React.createElement("span", {
    className: "stat-count"
  }, fmtUptime(metrics.system.uptime_seconds)), /*#__PURE__*/React.createElement("span", {
    className: "stat-label"
  }, "Uptime"))), /*#__PURE__*/React.createElement("h2", {
    className: "section-title"
  }, "Requests"), /*#__PURE__*/React.createElement("div", {
    className: "stats-grid"
  }, /*#__PURE__*/React.createElement("div", {
    className: "stat-card"
  }, /*#__PURE__*/React.createElement("span", {
    className: "stat-count"
  }, metrics.requests.total), /*#__PURE__*/React.createElement("span", {
    className: "stat-label"
  }, "Total")), /*#__PURE__*/React.createElement("div", {
    className: "stat-card"
  }, /*#__PURE__*/React.createElement("span", {
    className: "stat-count"
  }, metrics.requests.avg_latency_ms.toFixed(1), "ms"), /*#__PURE__*/React.createElement("span", {
    className: "stat-label"
  }, "Avg Latency"))), Object.keys(metrics.requests.by_endpoint).length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "table-wrap",
    style: {
      marginTop: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Endpoint"), /*#__PURE__*/React.createElement("th", null, "Count"), /*#__PURE__*/React.createElement("th", null, "Avg Latency"), /*#__PURE__*/React.createElement("th", null, "Status Codes"))), /*#__PURE__*/React.createElement("tbody", null, Object.entries(metrics.requests.by_endpoint).map(([path, ep]) => /*#__PURE__*/React.createElement("tr", {
    key: path
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, path)), /*#__PURE__*/React.createElement("td", null, ep.count), /*#__PURE__*/React.createElement("td", null, ep.avg_latency_ms.toFixed(1), "ms"), /*#__PURE__*/React.createElement("td", null, Object.entries(ep.status_count).map(([code, count]) => /*#__PURE__*/React.createElement(__ds_scope.Badge, {
    key: code,
    variant: code.startsWith('2') ? 'success' : code.startsWith('4') ? 'warning' : 'error',
    style: {
      marginRight: 'var(--space-2)'
    }
  }, code, ": ", count))))))))), tab === 'logs' && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("form", {
    onSubmit: e => {
      e.preventDefault();
      fetchLogs(logQuery);
    },
    className: "form-row"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Search logs...",
    value: logQuery,
    onChange: e => setLogQuery(e.target.value)
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    size: "sm"
  }, "Search"))), logs.length === 0 ? /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No logs.") : /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Level"), /*#__PURE__*/React.createElement("th", null, "Message"), /*#__PURE__*/React.createElement("th", null, "Time"))), /*#__PURE__*/React.createElement("tbody", null, logs.map(l => /*#__PURE__*/React.createElement("tr", {
    key: l.id
  }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Badge, {
    variant: l.level === 'error' ? 'error' : l.level === 'warn' ? 'warning' : 'success'
  }, l.level)), /*#__PURE__*/React.createElement("td", {
    className: "max"
  }, l.message), /*#__PURE__*/React.createElement("td", null, new Date(l.created_at).toLocaleString()))))))), tab === 'alerts' && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "primary",
    size: "sm",
    onClick: () => setShowAlertForm(!showAlertForm)
  }, showAlertForm ? 'Cancel' : 'New Alert')), showAlertForm && /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("form", {
    onSubmit: handleCreateAlert,
    className: "fn-form"
  }, /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Name *",
    value: alertForm.name,
    onChange: e => setAlertForm(p => ({
      ...p,
      name: e.target.value
    })),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Metric *",
    value: alertForm.metric,
    onChange: e => setAlertForm(p => ({
      ...p,
      metric: e.target.value
    })),
    required: true
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    placeholder: "Threshold",
    type: "number",
    step: "0.1",
    value: alertForm.threshold,
    onChange: e => setAlertForm(p => ({
      ...p,
      threshold: +e.target.value
    }))
  }), /*#__PURE__*/React.createElement(__ds_scope.Input, {
    as: "select",
    value: alertForm.operator,
    onChange: e => setAlertForm(p => ({
      ...p,
      operator: e.target.value
    }))
  }, /*#__PURE__*/React.createElement("option", {
    value: "gt"
  }, "Greater Than"), /*#__PURE__*/React.createElement("option", {
    value: "lt"
  }, "Less Than"), /*#__PURE__*/React.createElement("option", {
    value: "eq"
  }, "Equals")), /*#__PURE__*/React.createElement("label", {
    className: "checkbox-label"
  }, /*#__PURE__*/React.createElement("input", {
    type: "checkbox",
    checked: alertForm.enabled,
    onChange: e => setAlertForm(p => ({
      ...p,
      enabled: e.target.checked
    }))
  }), "Enabled"), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit"
  }, "Create"))), alerts.length === 0 ? /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No alerts.") : /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Metric"), /*#__PURE__*/React.createElement("th", null, "Threshold"), /*#__PURE__*/React.createElement("th", null, "Operator"), /*#__PURE__*/React.createElement("th", null, "Status"))), /*#__PURE__*/React.createElement("tbody", null, alerts.map(a => /*#__PURE__*/React.createElement("tr", {
    key: a.id
  }, /*#__PURE__*/React.createElement("td", null, a.name), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, a.metric)), /*#__PURE__*/React.createElement("td", null, a.threshold), /*#__PURE__*/React.createElement("td", null, a.operator), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Badge, {
    variant: a.enabled ? 'success' : 'warning'
  }, a.enabled ? 'Enabled' : 'Disabled')))))))));
}
Object.assign(__ds_scope, { MonitoringPage, __ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse: MonitoringPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/MonitoringPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/NotFoundPage.tsx
try { (() => {
function NotFoundPage() {
  return /*#__PURE__*/React.createElement("div", {
    className: "not-found"
  }, /*#__PURE__*/React.createElement("h1", null, "404"), /*#__PURE__*/React.createElement("p", null, "Page not found"), /*#__PURE__*/React.createElement(Link, {
    to: "/"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "primary",
    style: {
      marginTop: 'var(--space-8)'
    }
  }, "Go to Dashboard")));
}
Object.assign(__ds_scope, { NotFoundPage, __ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t: NotFoundPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/NotFoundPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/SqlEditorPage.tsx
try { (() => {
const {
  useState
} = React;
function SqlEditorPage() {
  const [query, setQuery] = useState("SELECT name FROM sqlite_master WHERE type = 'table' ORDER BY name");
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');
  const [running, setRunning] = useState(false);
  const handleRun = async () => {
    setError('');
    setResult(null);
    setRunning(true);
    try {
      const res = await fetch('/api/sql', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          query
        })
      });
      const data = await res.json();
      if (!res.ok) {
        setError(data.error || `error: ${res.status}`);
        return;
      }
      setResult(data);
    } catch {
      setError('network error');
    } finally {
      setRunning(false);
    }
  };
  const handleKeyDown = e => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') {
      e.preventDefault();
      handleRun();
    }
  };
  return /*#__PURE__*/React.createElement("div", {
    className: "sql-editor"
  }, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "SQL Editor"
  }), /*#__PURE__*/React.createElement("div", {
    className: "sql-layout"
  }, /*#__PURE__*/React.createElement("div", {
    className: "sql-input-area"
  }, /*#__PURE__*/React.createElement("textarea", {
    className: "sql-textarea",
    value: query,
    onChange: e => setQuery(e.target.value),
    onKeyDown: handleKeyDown,
    rows: 8,
    spellCheck: false,
    placeholder: "Enter SQL query..."
  }), /*#__PURE__*/React.createElement("div", {
    className: "sql-actions"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    onClick: handleRun,
    disabled: running || !query.trim()
  }, running ? 'Running...' : 'Run (⌘⏎)'))), error && /*#__PURE__*/React.createElement("p", {
    className: "sql-error"
  }, error), result && /*#__PURE__*/React.createElement("div", {
    className: "sql-results"
  }, /*#__PURE__*/React.createElement("p", {
    className: "result-meta"
  }, result.rows.length, " row", result.rows.length !== 1 ? 's' : '', " returned"), result.columns.length > 0 && result.rows.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, result.columns.map(col => /*#__PURE__*/React.createElement("th", {
    key: col
  }, col)))), /*#__PURE__*/React.createElement("tbody", null, result.rows.map((row, i) => /*#__PURE__*/React.createElement("tr", {
    key: i
  }, result.columns.map(col => /*#__PURE__*/React.createElement("td", {
    className: "max",
    key: col
  }, formatCell(row[col])))))))))));
}
function formatCell(val) {
  if (val === null || val === undefined) return 'NULL';
  if (typeof val === 'object') return JSON.stringify(val);
  return String(val);
}
Object.assign(__ds_scope, { SqlEditorPage, __ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr: SqlEditorPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/SqlEditorPage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/StoragePage.tsx
try { (() => {
const {
  useEffect,
  useState,
  useRef
} = React;
function fmtSize(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
function StoragePage() {
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [uploading, setUploading] = useState(false);
  const fileRef = useRef(null);
  const fetchFiles = async () => {
    setLoading(true);
    try {
      const res = await fetch('/api/storage/files');
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || `error: ${res.status}`);
        setFiles([]);
      } else {
        setFiles(d.data || []);
      }
    } catch {
      setError('network error');
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    fetchFiles();
  }, []);
  const handleUpload = async e => {
    e.preventDefault();
    const input = fileRef.current;
    if (!input || !input.files?.length) return;
    setUploading(true);
    setError('');
    const fd = new FormData();
    fd.append('file', input.files[0]);
    try {
      const res = await fetch('/api/storage/upload', {
        method: 'POST',
        body: fd
      });
      const d = await res.json();
      if (!res.ok) {
        setError(d.error || 'upload failed');
        setUploading(false);
        return;
      }
      input.value = '';
      fetchFiles();
    } catch {
      setError('network error');
    } finally {
      setUploading(false);
    }
  };
  const handleDelete = async id => {
    if (!confirm('Delete this file?')) return;
    try {
      const res = await fetch(`/api/storage/files/${id}`, {
        method: 'DELETE'
      });
      if (!res.ok) {
        const d = await res.json();
        setError(d.error || 'delete failed');
        return;
      }
      setFiles(prev => prev.filter(f => f.id !== id));
    } catch {
      setError('network error');
    }
  };
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading files...");
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Storage"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: fetchFiles
  }, "Refresh")), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), /*#__PURE__*/React.createElement("div", {
    className: "card",
    style: {
      marginBottom: 'var(--space-8)'
    }
  }, /*#__PURE__*/React.createElement("form", {
    onSubmit: handleUpload,
    className: "upload-form"
  }, /*#__PURE__*/React.createElement("input", {
    type: "file",
    ref: fileRef,
    required: true,
    style: {
      fontSize: 'var(--text-s)'
    }
  }), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    type: "submit",
    size: "sm",
    disabled: uploading
  }, uploading ? 'Uploading...' : 'Upload'))), files.length === 0 && !error && /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No files uploaded."), files.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Name"), /*#__PURE__*/React.createElement("th", null, "Size"), /*#__PURE__*/React.createElement("th", null, "Type"), /*#__PURE__*/React.createElement("th", null, "Uploaded"), /*#__PURE__*/React.createElement("th", null, "Actions"))), /*#__PURE__*/React.createElement("tbody", null, files.map(f => /*#__PURE__*/React.createElement("tr", {
    key: f.id
  }, /*#__PURE__*/React.createElement("td", null, f.name), /*#__PURE__*/React.createElement("td", null, fmtSize(f.size)), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("code", null, f.mime_type)), /*#__PURE__*/React.createElement("td", null, new Date(f.created_at).toLocaleString()), /*#__PURE__*/React.createElement("td", {
    className: "actions-cell"
  }, /*#__PURE__*/React.createElement("a", {
    href: `/api/storage/files/${f.id}`,
    className: "btn btn-secondary btn-sm",
    download: f.name,
    style: {
      marginRight: 'var(--space-2)'
    }
  }, "Download"), /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "danger",
    size: "sm",
    onClick: () => handleDelete(f.id)
  }, "Delete"))))))));
}
Object.assign(__ds_scope, { StoragePage, __ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p: StoragePage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/StoragePage.tsx", error: String((e && e.message) || e) }); }

// source/bigbase-ui/pages/UsersPage.tsx
try { (() => {
const {
  useEffect,
  useState
} = React;
function UsersPage() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const fetchUsers = async () => {
    setError('');
    setLoading(true);
    try {
      const res = await fetch('/api/auth/users');
      const data = await res.json();
      if (!res.ok) {
        setError(data.error || `error: ${res.status}`);
        setUsers([]);
        return;
      }
      setUsers(data.data || []);
    } catch {
      setError('network error');
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    fetchUsers();
  }, []);
  const handleDelete = async id => {
    if (!confirm(`Delete user #${id}?`)) return;
    try {
      const res = await fetch(`/api/auth/users/${id}`, {
        method: 'DELETE'
      });
      if (!res.ok) {
        const data = await res.json();
        setError(data.error || `error: ${res.status}`);
        return;
      }
      setUsers(prev => prev.filter(u => u.id !== id));
    } catch {
      setError('network error');
    }
  };
  if (loading) return /*#__PURE__*/React.createElement("div", {
    className: "loading"
  }, "Loading users...");
  return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(__ds_scope.PageHeader, {
    title: "Users"
  }, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "secondary",
    size: "sm",
    onClick: fetchUsers
  }, "Refresh")), error && /*#__PURE__*/React.createElement("p", {
    className: "input-error-text"
  }, error), users.length === 0 && !error && /*#__PURE__*/React.createElement("p", {
    className: "dim"
  }, "No users found."), users.length > 0 && /*#__PURE__*/React.createElement("div", {
    className: "table-wrap"
  }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "ID"), /*#__PURE__*/React.createElement("th", null, "Email"), /*#__PURE__*/React.createElement("th", null, "Created"), /*#__PURE__*/React.createElement("th", null, "Actions"))), /*#__PURE__*/React.createElement("tbody", null, users.map(u => /*#__PURE__*/React.createElement("tr", {
    key: u.id
  }, /*#__PURE__*/React.createElement("td", null, u.id), /*#__PURE__*/React.createElement("td", null, u.email), /*#__PURE__*/React.createElement("td", null, new Date(u.created_at).toLocaleString()), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(__ds_scope.Button, {
    variant: "danger",
    size: "sm",
    onClick: () => handleDelete(u.id)
  }, "Delete"))))))));
}
Object.assign(__ds_scope, { UsersPage, __ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui: UsersPage });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/pages/UsersPage.tsx", error: String((e && e.message) || e) }); }

__ds_scope.__ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1$13xveek = __ds_scope.__ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1 !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1 : __ds_scope.LoginPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_DashboardPage_120zej4$mogxob = __ds_scope.__ds_default_source_bigbase_ui_pages_DashboardPage_120zej4 !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_DashboardPage_120zej4 : __ds_scope.DashboardPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6$txrdad = __ds_scope.__ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6 !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6 : __ds_scope.DataStudioPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr$1reqlu2 = __ds_scope.__ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr : __ds_scope.SqlEditorPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui$1n11yp1 = __ds_scope.__ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui : __ds_scope.UsersPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5$15vgx4w = __ds_scope.__ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5 !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5 : __ds_scope.GitReposPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_DeployPage_17gbuph$7njolc = __ds_scope.__ds_default_source_bigbase_ui_pages_DeployPage_17gbuph !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_DeployPage_17gbuph : __ds_scope.DeployPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza$1579x4h = __ds_scope.__ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza : __ds_scope.MessagingPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p$mt7puw = __ds_scope.__ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p : __ds_scope.StoragePage;

__ds_scope.__ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx$1socyh8 = __ds_scope.__ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx : __ds_scope.FunctionsPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf$8ghtly = __ds_scope.__ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf : __ds_scope.ForgePage;

__ds_scope.__ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c$n2l73f = __ds_scope.__ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c : __ds_scope.CiciPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse$1r6cjcp = __ds_scope.__ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse : __ds_scope.MonitoringPage;

__ds_scope.__ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t$m6vebk = __ds_scope.__ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t : __ds_scope.NotFoundPage;

__ds_scope.__ds_default_source_bigbase_ui_Layout_2g8xmy$1x3lec = __ds_scope.__ds_default_source_bigbase_ui_Layout_2g8xmy !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_Layout_2g8xmy : __ds_scope.Layout;

// source/bigbase-ui/App.tsx
try { (() => {
function App() {
  return /*#__PURE__*/React.createElement(Routes, null, /*#__PURE__*/React.createElement(Route, {
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_Layout_2g8xmy$1x3lec, null)
  }, /*#__PURE__*/React.createElement(Route, {
    path: "/",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_DashboardPage_120zej4$mogxob, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "data",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6$txrdad, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "sql",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr$1reqlu2, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "users",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui$1n11yp1, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "repos",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5$15vgx4w, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "deploy",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_DeployPage_17gbuph$7njolc, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "messaging",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza$1579x4h, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "storage",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p$mt7puw, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "functions",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx$1socyh8, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "forge",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf$8ghtly, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "cici",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c$n2l73f, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "monitoring",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse$1r6cjcp, null)
  }), /*#__PURE__*/React.createElement(Route, {
    path: "*",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t$m6vebk, null)
  })), /*#__PURE__*/React.createElement(Route, {
    path: "login",
    element: /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1$13xveek, null)
  }));
}
Object.assign(__ds_scope, { App, __ds_default_source_bigbase_ui_App_tku48d: App });
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/App.tsx", error: String((e && e.message) || e) }); }

__ds_scope.__ds_default_source_bigbase_ui_App_tku48d$1urd5a3 = __ds_scope.__ds_default_source_bigbase_ui_App_tku48d !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_App_tku48d : __ds_scope.App;

// source/bigbase-ui/main.tsx
try { (() => {
const {
  StrictMode
} = React;
createRoot(document.getElementById('root')).render(/*#__PURE__*/React.createElement(StrictMode, null, /*#__PURE__*/React.createElement(HashRouter, null, /*#__PURE__*/React.createElement(__ds_scope.__ds_default_source_bigbase_ui_App_tku48d$1urd5a3, null))));
})(); } catch (e) { __ds_ns.__errors.push({ path: "source/bigbase-ui/main.tsx", error: String((e && e.message) || e) }); }

// ui_kits/admin-console/app.jsx
try { (() => {
/* BigBase Admin Console — app shell */
(function () {
  const {
    useState,
    useEffect
  } = React;
  const {
    Sidebar,
    ToastProvider,
    useToast,
    PreviewBanner,
    Login,
    Dashboard,
    SitesList,
    SiteDetail,
    CreateSite,
    Functions,
    FunctionDetail,
    DataStudio,
    Messaging,
    MessagingDetail,
    Settings,
    Placeholder
  } = window;
  function Shell() {
    const [user, setUser] = useState(null);
    const [route, setRoute] = useState('dashboard');
    const [dark, setDark] = useState(false);
    const [theme, setTheme] = useState(() => {
      try {
        return localStorage.getItem('bigbase-theme') || 'default';
      } catch (e) {
        return 'default';
      }
    });
    const [sites, setSites] = useState(window.DATA.SITES.map(s => ({
      ...s
    })));
    const [loadingSites, setLoadingSites] = useState(false);
    const toast = useToast();
    useEffect(() => {
      document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
    }, [dark]);
    useEffect(() => {
      const root = document.documentElement;
      if (theme && theme !== 'default') root.setAttribute('data-accent', theme);else root.removeAttribute('data-accent');
      try {
        localStorage.setItem('bigbase-theme', theme);
      } catch (e) {}
    }, [theme]);
    const nav = r => {
      if (r === 'sites') {
        setLoadingSites(true);
        setTimeout(() => setLoadingSites(false), 550);
      }
      setRoute(r);
      window.scrollTo(0, 0);
    };
    const onDeployed = site => {
      setSites(prev => [site, ...prev.filter(s => s.id !== site.id)]);
      toast({
        type: 'success',
        title: 'Deployment ready',
        msg: `${site.name} is live at ${site.url}`
      });
      nav('sites/' + site.id);
    };
    const onRedeploy = site => {
      toast({
        type: 'info',
        title: 'Redeploy started',
        msg: `Building ${site.name} from ${site.branch}…`
      });
      setSites(prev => prev.map(s => s.id === site.id ? {
        ...s,
        status: 'building'
      } : s));
      setTimeout(() => {
        setSites(prev => prev.map(s => s.id === site.id ? {
          ...s,
          status: 'ready',
          updated: Date.now()
        } : s));
        toast({
          type: 'success',
          title: 'Deployment ready',
          msg: `${site.name} redeployed successfully.`
        });
      }, 2600);
    };
    if (!user) return /*#__PURE__*/React.createElement(Login, {
      onLogin: email => {
        setUser({
          email
        });
        nav('dashboard');
      }
    });
    let screen;
    if (route === 'dashboard') screen = /*#__PURE__*/React.createElement(Dashboard, {
      user: user,
      onNav: nav
    });else if (route === 'sites') screen = /*#__PURE__*/React.createElement(SitesList, {
      sites: sites,
      onNav: nav,
      loading: loadingSites
    });else if (route === 'sites/new') screen = /*#__PURE__*/React.createElement(CreateSite, {
      onCancel: () => nav('sites'),
      onBackToList: () => nav('sites'),
      onDeployed: onDeployed
    });else if (route.startsWith('sites/')) {
      const site = sites.find(s => s.id === route.slice(6));
      screen = site ? /*#__PURE__*/React.createElement(SiteDetail, {
        site: site,
        onNav: nav,
        onRedeploy: onRedeploy
      }) : /*#__PURE__*/React.createElement(Placeholder, {
        title: "Site not found"
      });
    } else if (route === 'functions') screen = /*#__PURE__*/React.createElement(Functions, {
      onNav: nav
    });else if (route.startsWith('functions/')) {
      const fn = window.DATA.FUNCTIONS.find(f => f.id === route.slice(10));
      screen = fn ? /*#__PURE__*/React.createElement(FunctionDetail, {
        fn: fn,
        onNav: nav
      }) : /*#__PURE__*/React.createElement(Placeholder, {
        title: "Function not found"
      });
    } else if (route === 'data') screen = /*#__PURE__*/React.createElement(DataStudio, null);else if (route === 'messaging') screen = /*#__PURE__*/React.createElement(Messaging, {
      onNav: nav
    });else if (route.startsWith('messaging/')) {
      const tpl = window.DATA.TEMPLATES.find(t => t.id === route.slice(10));
      screen = tpl ? /*#__PURE__*/React.createElement(MessagingDetail, {
        tpl: tpl,
        onNav: nav
      }) : /*#__PURE__*/React.createElement(Placeholder, {
        title: "Template not found"
      });
    } else if (route === 'settings') screen = /*#__PURE__*/React.createElement(Settings, {
      user: user
    });else {
      const titles = {
        sql: 'SQL Editor',
        storage: 'Storage',
        users: 'Users',
        repos: 'Git Repos',
        cici: 'CI / CD',
        monitoring: 'Monitoring'
      };
      screen = /*#__PURE__*/React.createElement(Placeholder, {
        title: titles[route] || 'BigBase'
      });
    }
    return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(PreviewBanner, null), /*#__PURE__*/React.createElement("div", {
      className: "layout"
    }, /*#__PURE__*/React.createElement(Sidebar, {
      route: route,
      onNav: nav,
      user: user,
      dark: dark,
      onToggleDark: () => setDark(d => !d),
      onLogout: () => setUser(null),
      theme: theme,
      onTheme: setTheme
    }), /*#__PURE__*/React.createElement("main", {
      className: "content"
    }, screen)));
  }
  function App() {
    return /*#__PURE__*/React.createElement(ToastProvider, null, /*#__PURE__*/React.createElement(Shell, null));
  }
  ReactDOM.createRoot(document.getElementById('root')).render(/*#__PURE__*/React.createElement(App, null));
})();
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/admin-console/app.jsx", error: String((e && e.message) || e) }); }

// ui_kits/admin-console/data.jsx
try { (() => {
/* BigBase Admin Console — mock data + helpers */
(function () {
  const now = Date.now();
  const min = 60000,
    hr = 3600000,
    day = 86400000;
  const SITES = [{
    id: 'ste_marketing',
    name: 'marketing-site',
    framework: 'Astro',
    repo: 'danielvm/marketing',
    branch: 'main',
    root: '/',
    status: 'ready',
    url: 'marketing-site.bigbase.local',
    commit: 'a3f91c2',
    commitMsg: 'Update hero copy and pricing grid',
    updated: now - 2 * min,
    env: 'production',
    thumb: 'grad-indigo',
    deployments: [{
      id: 'dep_1',
      status: 'ready',
      branch: 'main',
      commit: 'a3f91c2',
      msg: 'Update hero copy and pricing grid',
      when: now - 2 * min,
      duration: '38s'
    }, {
      id: 'dep_2',
      status: 'ready',
      branch: 'main',
      commit: '7e10b4d',
      msg: 'Add testimonials section',
      when: now - 6 * hr,
      duration: '41s'
    }, {
      id: 'dep_3',
      status: 'ready',
      branch: 'main',
      commit: '2c8aa90',
      msg: 'Initial import',
      when: now - 3 * day,
      duration: '52s'
    }]
  }, {
    id: 'ste_docs',
    name: 'docs',
    framework: 'Next.js',
    repo: 'danielvm/docs',
    branch: 'main',
    root: '/',
    status: 'building',
    url: 'docs.bigbase.local',
    commit: 'f29d1aa',
    commitMsg: 'Restructure API reference nav',
    updated: now - 40000,
    env: 'production',
    thumb: 'grad-slate',
    progress: 62,
    deployments: [{
      id: 'dep_4',
      status: 'building',
      branch: 'main',
      commit: 'f29d1aa',
      msg: 'Restructure API reference nav',
      when: now - 40000,
      duration: '—'
    }, {
      id: 'dep_5',
      status: 'ready',
      branch: 'main',
      commit: 'b81f0c3',
      msg: 'Fix broken anchor links',
      when: now - day,
      duration: '1m 12s'
    }]
  }, {
    id: 'ste_app',
    name: 'dashboard-app',
    framework: 'Vite + React',
    repo: 'danielvm/dashboard',
    branch: 'release',
    root: '/web',
    status: 'failed',
    url: 'dashboard-app.bigbase.local',
    commit: '0d4e7f1',
    commitMsg: 'Bump dependencies for security patch',
    updated: now - 25 * min,
    env: 'preview',
    thumb: 'grad-rose',
    deployments: [{
      id: 'dep_6',
      status: 'failed',
      branch: 'release',
      commit: '0d4e7f1',
      msg: 'Bump dependencies for security patch',
      when: now - 25 * min,
      duration: '19s'
    }, {
      id: 'dep_7',
      status: 'ready',
      branch: 'release',
      commit: 'ac5512e',
      msg: 'Wire up auth guard',
      when: now - 2 * day,
      duration: '1m 04s'
    }]
  }];
  const REPOS = [{
    name: 'danielvm/marketing',
    desc: 'Astro marketing site',
    updated: '2 hours ago',
    lang: 'Astro',
    private: false
  }, {
    name: 'danielvm/dashboard',
    desc: 'Internal admin dashboard',
    updated: 'yesterday',
    lang: 'TypeScript',
    private: true
  }, {
    name: 'danielvm/api-gateway',
    desc: 'Go edge gateway + proxy',
    updated: '3 days ago',
    lang: 'Go',
    private: true
  }, {
    name: 'danielvm/blog-engine',
    desc: 'Markdown blog with RSS',
    updated: 'last week',
    lang: 'JavaScript',
    private: false
  }, {
    name: 'danielvm/landing-kit',
    desc: 'Reusable landing components',
    updated: '2 weeks ago',
    lang: 'Svelte',
    private: false
  }];
  const FRAMEWORKS = [{
    id: 'astro',
    name: 'Astro',
    build: 'npm run build',
    output: 'dist',
    install: 'npm install'
  }, {
    id: 'next',
    name: 'Next.js',
    build: 'next build',
    output: '.next',
    install: 'npm install'
  }, {
    id: 'vite',
    name: 'Vite',
    build: 'vite build',
    output: 'dist',
    install: 'npm install'
  }, {
    id: 'static',
    name: 'Static HTML',
    build: '—',
    output: '.',
    install: '—'
  }, {
    id: 'sveltekit',
    name: 'SvelteKit',
    build: 'vite build',
    output: 'build',
    install: 'npm install'
  }];
  const FUNCTIONS = [{
    id: 'fn_resize',
    name: 'image-resize',
    runtime: 'JavaScript',
    trigger: 'HTTP',
    timeout: 30,
    updated: now - 3 * hr,
    status: 'active',
    created: 'Feb 14, 2026',
    deployed: now - 3 * hr,
    invocations: '12.4k',
    errors: '0.2%',
    env: [{
      k: 'MAX_WIDTH',
      v: '2048'
    }, {
      k: 'CDN_BUCKET',
      v: 'media-prod'
    }],
    code: "export default async function (req, res) {\n  const { url, width } = req.query;\n  const img = await fetch(url).then(r => r.arrayBuffer());\n  const out = await resize(img, { width: Number(width) });\n  res.setHeader('content-type', 'image/webp');\n  return res.send(out);\n}",
    logs: [{
      t: '12:04:21',
      m: 'GET /image-resize?width=640 → 200 (84ms)',
      k: 'ok'
    }, {
      t: '12:03:58',
      m: 'GET /image-resize?width=1280 → 200 (102ms)',
      k: 'ok'
    }, {
      t: '12:01:12',
      m: 'cold start: 412ms',
      k: 'warn'
    }]
  }, {
    id: 'fn_webhook',
    name: 'stripe-webhook',
    runtime: 'JavaScript',
    trigger: 'HTTP',
    timeout: 15,
    updated: now - day,
    status: 'active',
    created: 'Jan 30, 2026',
    deployed: now - day,
    invocations: '3.1k',
    errors: '0.0%',
    env: [{
      k: 'STRIPE_SECRET',
      v: '••••••••'
    }],
    code: "export default async function (req, res) {\n  const event = verify(req.body, req.headers['stripe-signature']);\n  if (event.type === 'invoice.paid') await grantAccess(event.data);\n  return res.json({ received: true });\n}",
    logs: [{
      t: '09:18:02',
      m: 'POST /stripe-webhook → 200 (47ms)',
      k: 'ok'
    }, {
      t: '08:54:30',
      m: 'invoice.paid handled for cus_J2x',
      k: 'ok'
    }]
  }, {
    id: 'fn_cron',
    name: 'nightly-digest',
    runtime: 'JavaScript',
    trigger: 'Schedule',
    timeout: 60,
    updated: now - 2 * day,
    status: 'active',
    created: 'Dec 2, 2025',
    deployed: now - 2 * day,
    invocations: '186',
    errors: '1.1%',
    schedule: '0 6 * * *',
    env: [{
      k: 'DIGEST_FROM',
      v: 'hello@bigbase.dev'
    }],
    code: "export default async function () {\n  const users = await db.users.where({ digest: true });\n  for (const u of users) await sendDigest(u);\n}",
    logs: [{
      t: '06:00:03',
      m: 'cron triggered (0 6 * * *)',
      k: ''
    }, {
      t: '06:00:41',
      m: 'sent 142 digests in 38s',
      k: 'ok'
    }]
  }];
  const RUNTIMES = ['JavaScript (Node 20)', 'TypeScript', 'Python 3.12', 'Go 1.22', 'Bun'];
  const TRIGGERS = ['HTTP', 'Schedule', 'Event'];

  /* ── Messaging templates ── */
  const TEMPLATES = [{
    id: 'tpl_welcome',
    name: 'Welcome email',
    type: 'Email',
    status: 'active',
    updated: now - 4 * hr,
    sends: '2.1k',
    subject: 'Welcome to {{workspace}} 🎉',
    body: 'Hi {{name}},\n\nThanks for joining {{workspace}}. Your account is ready — deploy your first site in minutes.\n\nIf you have questions, just reply to this email.\n\n— The {{workspace}} team',
    vars: ['name', 'workspace']
  }, {
    id: 'tpl_reset',
    name: 'Password reset',
    type: 'Email',
    status: 'active',
    updated: now - 2 * day,
    sends: '540',
    subject: 'Reset your {{workspace}} password',
    body: 'Hi {{name}},\n\nWe got a request to reset your password. Use the link below within 30 minutes:\n\n{{reset_url}}\n\nDidn\'t request this? You can safely ignore this email.',
    vars: ['name', 'workspace', 'reset_url']
  }, {
    id: 'tpl_deploy',
    name: 'Deploy succeeded',
    type: 'Email',
    status: 'active',
    updated: now - 6 * hr,
    sends: '8.7k',
    subject: '{{site}} is live',
    body: 'Your deployment of {{site}} finished in {{duration}}.\n\nView it at {{url}}.',
    vars: ['site', 'duration', 'url']
  }, {
    id: 'tpl_otp',
    name: 'Login code',
    type: 'SMS',
    status: 'draft',
    updated: now - day,
    sends: '—',
    subject: '—',
    body: 'Your {{workspace}} login code is {{code}}. It expires in 10 minutes.',
    vars: ['workspace', 'code']
  }, {
    id: 'tpl_invite',
    name: 'Team invitation',
    type: 'Email',
    status: 'paused',
    updated: now - 3 * day,
    sends: '312',
    subject: '{{inviter}} invited you to {{workspace}}',
    body: '{{inviter}} added you to the {{workspace}} workspace on BigBase.\n\nAccept the invite to get started:\n\n{{accept_url}}',
    vars: ['inviter', 'workspace', 'accept_url']
  }];
  const TEMPLATE_TYPES = ['Email', 'SMS', 'Push'];
  const COLLECTIONS = ['users', 'sites', 'deployments', 'sessions', 'audit_log'];
  /* schema for the active collection (Data Studio schema explorer) */
  const SCHEMA = {
    users: [{
      name: 'id',
      type: 'string',
      pk: true,
      nullable: false,
      def: 'uuid()'
    }, {
      name: 'email',
      type: 'string',
      pk: false,
      nullable: false,
      def: '—'
    }, {
      name: 'role',
      type: 'enum',
      pk: false,
      nullable: false,
      def: "'member'"
    }, {
      name: 'verified',
      type: 'boolean',
      pk: false,
      nullable: false,
      def: 'false'
    }, {
      name: 'created',
      type: 'datetime',
      pk: false,
      nullable: false,
      def: 'now()'
    }]
  };
  const TABLE_ROWS = [{
    id: 'usr_8x2a',
    email: 'daniel@bigbase.dev',
    role: 'owner',
    verified: true,
    created: 'Mar 2, 2026'
  }, {
    id: 'usr_91kd',
    email: 'maya@acme.io',
    role: 'admin',
    verified: true,
    created: 'Apr 11, 2026'
  }, {
    id: 'usr_3jf0',
    email: 'sam@studio.co',
    role: 'member',
    verified: false,
    created: 'May 7, 2026'
  }, {
    id: 'usr_55ze',
    email: 'lee@northwind.com',
    role: 'member',
    verified: true,
    created: 'May 22, 2026'
  }];
  const BUILD_LOG = [{
    t: '00:00',
    m: 'Cloning danielvm/marketing @ main…',
    k: ''
  }, {
    t: '00:02',
    m: 'Checked out commit a3f91c2',
    k: 'ok'
  }, {
    t: '00:03',
    m: 'Detected framework: Astro',
    k: 'ok'
  }, {
    t: '00:04',
    m: 'Running: npm install',
    k: ''
  }, {
    t: '00:21',
    m: 'added 412 packages in 16s',
    k: ''
  }, {
    t: '00:22',
    m: 'Running: npm run build',
    k: ''
  }, {
    t: '00:34',
    m: 'building client + server…',
    k: ''
  }, {
    t: '00:37',
    m: 'Completed in 14.82s — 18 pages built',
    k: 'ok'
  }, {
    t: '00:38',
    m: 'Uploading dist → edge',
    k: ''
  }, {
    t: '00:38',
    m: 'Deployment ready ✓  https://marketing-site.bigbase.local',
    k: 'ok'
  }];
  function timeAgo(ts) {
    const d = Date.now() - ts;
    if (d < min) return 'just now';
    if (d < hr) return Math.floor(d / min) + 'm ago';
    if (d < day) return Math.floor(d / hr) + 'h ago';
    return Math.floor(d / day) + 'd ago';
  }
  function statusVariant(s) {
    if (s === 'ready') return 'success';
    if (s === 'building' || s === 'pending') return 'warning';
    if (s === 'failed') return 'error';
    return 'neutral';
  }
  const THUMBS = {
    'grad-indigo': 'linear-gradient(135deg,#4F46E5 0%,#7C73F0 100%)',
    'grad-slate': 'linear-gradient(135deg,#3a3a48 0%,#5b5b6e 100%)',
    'grad-rose': 'linear-gradient(135deg,#9f1239 0%,#e11d48 100%)',
    'grad-emerald': 'linear-gradient(135deg,#065f46 0%,#10b981 100%)'
  };

  /* ── Month themes (accent presets). swatch = CSS for the dropdown dot. ── */
  const THEMES = [{
    key: 'default',
    month: 'Default',
    label: 'Indigo',
    swatch: 'rgb(79, 70, 229)'
  }, {
    key: 'january',
    month: 'January',
    label: 'Teal',
    swatch: 'rgb(13, 148, 136)'
  }, {
    key: 'february',
    month: 'February',
    label: 'Orange',
    swatch: 'rgb(234, 88, 12)'
  }, {
    key: 'march',
    month: 'March',
    label: 'Purple',
    swatch: 'rgb(124, 58, 237)'
  }, {
    key: 'april',
    month: 'April',
    label: 'Green',
    swatch: 'rgb(22, 163, 74)'
  }, {
    key: 'may',
    month: 'May',
    label: 'Lavender',
    swatch: 'rgb(167, 139, 250)'
  }, {
    key: 'june',
    month: 'June',
    label: 'Rainbow',
    swatch: 'linear-gradient(to right, rgb(239,68,68), rgb(245,158,11), rgb(16,185,129), rgb(59,130,246), rgb(139,92,246))'
  }, {
    key: 'july',
    month: 'July',
    label: 'Peach',
    swatch: 'rgb(253, 186, 116)'
  }, {
    key: 'august',
    month: 'August',
    label: 'Silver',
    swatch: 'rgb(156, 163, 175)'
  }, {
    key: 'september',
    month: 'September',
    label: 'Yellow',
    swatch: 'rgb(234, 179, 8)'
  }, {
    key: 'october',
    month: 'October',
    label: 'Pink',
    swatch: 'rgb(236, 72, 153)'
  }, {
    key: 'november',
    month: 'November',
    label: 'Blue',
    swatch: 'rgb(37, 99, 235)'
  }, {
    key: 'december',
    month: 'December',
    label: 'Red',
    swatch: 'rgb(220, 38, 38)'
  }];
  window.DATA = {
    SITES,
    REPOS,
    FRAMEWORKS,
    FUNCTIONS,
    RUNTIMES,
    TRIGGERS,
    TEMPLATES,
    TEMPLATE_TYPES,
    COLLECTIONS,
    SCHEMA,
    TABLE_ROWS,
    BUILD_LOG,
    THEMES,
    timeAgo,
    statusVariant,
    THUMBS
  };
})();
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/admin-console/data.jsx", error: String((e && e.message) || e) }); }

// ui_kits/admin-console/screens.jsx
try { (() => {
/* BigBase Admin Console — core screens */
(function () {
  const {
    useState,
    useEffect
  } = React;
  const Icon = window.Icon;
  const {
    SITES,
    FUNCTIONS,
    COLLECTIONS,
    TABLE_ROWS,
    BUILD_LOG,
    timeAgo,
    statusVariant,
    THUMBS
  } = window.DATA;
  const {
    Badge,
    StatusBadge,
    Avatar,
    SkeletonCard
  } = window;
  function PageHeader({
    title,
    subtitle,
    children
  }) {
    return /*#__PURE__*/React.createElement("div", {
      className: "page-header"
    }, /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("h1", {
      className: "page-title"
    }, title), subtitle && /*#__PURE__*/React.createElement("div", {
      className: "page-subtitle"
    }, subtitle)), children && /*#__PURE__*/React.createElement("div", {
      className: "page-header-actions"
    }, children));
  }

  /* ── Login (refined: validation + password reset) ── */
  function Login({
    onLogin
  }) {
    const [mode, setMode] = useState('signin'); // signin | reset | reset-sent
    const [email, setEmail] = useState('daniel@bigbase.dev');
    const [pw, setPw] = useState('demo1234');
    const [errs, setErrs] = useState({});
    const [loading, setLoading] = useState(false);
    const submit = e => {
      e.preventDefault();
      const next = {};
      if (!email.includes('@')) next.email = 'Enter a valid email address.';
      if (pw.length < 6) next.pw = 'Password must be at least 6 characters.';
      setErrs(next);
      if (Object.keys(next).length) return;
      setLoading(true);
      setTimeout(() => {
        setLoading(false);
        onLogin(email);
      }, 700);
    };
    const sendReset = e => {
      e.preventDefault();
      if (!email.includes('@')) {
        setErrs({
          email: 'Enter a valid email address.'
        });
        return;
      }
      setErrs({});
      setLoading(true);
      setTimeout(() => {
        setLoading(false);
        setMode('reset-sent');
      }, 700);
    };
    return /*#__PURE__*/React.createElement("div", {
      className: "login-page",
      style: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        padding: 24,
        background: 'var(--bg-default)'
      }
    }, /*#__PURE__*/React.createElement("div", {
      style: {
        width: '100%',
        maxWidth: 400
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card",
      style: {
        padding: 40,
        borderRadius: 'var(--radius-l)',
        boxShadow: 'var(--shadow-m)'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        alignItems: 'center',
        gap: 8,
        marginBottom: 28
      }
    }, /*#__PURE__*/React.createElement("img", {
      src: "../../assets/bigbase-logo.svg",
      width: "44",
      height: "44",
      alt: "",
      style: {
        marginBottom: 6
      }
    }), /*#__PURE__*/React.createElement("div", {
      className: "ds-h2"
    }, mode === 'signin' ? 'Welcome back' : 'Reset your password'), /*#__PURE__*/React.createElement("div", {
      className: "ds-body dim",
      style: {
        textAlign: 'center'
      }
    }, mode === 'signin' && 'Sign in to your BigBase console', mode === 'reset' && 'Enter your email and we\'ll send a reset link.', mode === 'reset-sent' && `If an account exists for ${email}, a reset link is on its way.`)), mode === 'reset-sent' ? /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary btn-block",
      onClick: () => {
        setMode('signin');
        setErrs({});
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "arrow-left",
      size: 15
    }), " Back to sign in") : mode === 'reset' ? /*#__PURE__*/React.createElement("form", {
      className: "col gap-4",
      onSubmit: sendReset
    }, /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label",
      htmlFor: "lg-email"
    }, "Email"), /*#__PURE__*/React.createElement("input", {
      id: "lg-email",
      className: `input ${errs.email ? 'input-error' : ''}`,
      type: "email",
      value: email,
      onChange: e => setEmail(e.target.value),
      autoFocus: true
    }), errs.email && /*#__PURE__*/React.createElement("div", {
      className: "input-error-text"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "alert-triangle",
      size: 13,
      style: {
        marginRight: 4,
        verticalAlign: '-2px'
      }
    }), errs.email)), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary btn-block",
      type: "submit",
      disabled: loading
    }, loading ? /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("span", {
      className: "spinner spinner-sm"
    }), " Sending\u2026") : 'Send reset link'), /*#__PURE__*/React.createElement("button", {
      type: "button",
      className: "btn btn-link",
      style: {
        alignSelf: 'center'
      },
      onClick: () => {
        setMode('signin');
        setErrs({});
      }
    }, "Back to sign in")) : /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("form", {
      className: "col gap-4",
      onSubmit: submit,
      noValidate: true
    }, /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label",
      htmlFor: "lg-email"
    }, "Email"), /*#__PURE__*/React.createElement("input", {
      id: "lg-email",
      className: `input ${errs.email ? 'input-error' : ''}`,
      type: "email",
      value: email,
      onChange: e => setEmail(e.target.value),
      autoFocus: true
    }), errs.email && /*#__PURE__*/React.createElement("div", {
      className: "input-error-text"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "alert-triangle",
      size: 13,
      style: {
        marginRight: 4,
        verticalAlign: '-2px'
      }
    }), errs.email)), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label",
      htmlFor: "lg-pw"
    }, "Password"), /*#__PURE__*/React.createElement("button", {
      type: "button",
      className: "btn btn-link",
      style: {
        fontSize: 13
      },
      onClick: () => {
        setMode('reset');
        setErrs({});
      }
    }, "Forgot password?")), /*#__PURE__*/React.createElement("input", {
      id: "lg-pw",
      className: `input ${errs.pw ? 'input-error' : ''}`,
      type: "password",
      value: pw,
      onChange: e => setPw(e.target.value)
    }), errs.pw && /*#__PURE__*/React.createElement("div", {
      className: "input-error-text"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "alert-triangle",
      size: 13,
      style: {
        marginRight: 4,
        verticalAlign: '-2px'
      }
    }), errs.pw)), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary btn-block",
      type: "submit",
      disabled: loading
    }, loading ? /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("span", {
      className: "spinner spinner-sm"
    }), " Signing in\u2026") : 'Sign in')), /*#__PURE__*/React.createElement("div", {
      className: "divider"
    }, /*#__PURE__*/React.createElement("span", null, "or")), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary btn-block",
      onClick: () => onLogin(email)
    }, /*#__PURE__*/React.createElement("svg", {
      viewBox: "0 0 24 24",
      width: "16",
      height: "16"
    }, /*#__PURE__*/React.createElement("path", {
      fill: "#4285F4",
      d: "M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"
    }), /*#__PURE__*/React.createElement("path", {
      fill: "#34A853",
      d: "M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
    }), /*#__PURE__*/React.createElement("path", {
      fill: "#FBBC05",
      d: "M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
    }), /*#__PURE__*/React.createElement("path", {
      fill: "#EA4335",
      d: "M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
    })), "Continue with Google"))), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption",
      style: {
        textAlign: 'center',
        marginTop: 16
      }
    }, "Self-hosted \xB7 v2.0 \xB7 running on this instance")));
  }

  /* ── Dashboard hub ── */
  function Dashboard({
    user,
    onNav
  }) {
    const stats = [{
      label: 'Sites',
      count: SITES.length,
      icon: 'rocket',
      to: 'sites'
    }, {
      label: 'Functions',
      count: FUNCTIONS.length,
      icon: 'box',
      to: 'functions'
    }, {
      label: 'Git Repos',
      count: 5,
      icon: 'git-branch',
      to: 'repos'
    }, {
      label: 'Users',
      count: TABLE_ROWS.length,
      icon: 'users',
      to: 'users'
    }];
    const recent = SITES.flatMap(s => s.deployments.map(d => ({
      ...d,
      site: s.name
    }))).sort((a, b) => b.when - a.when).slice(0, 5);
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement(PageHeader, {
      title: `Welcome back, ${user.email.split('@')[0]}`,
      subtitle: "Here's what's running on your BigBase instance."
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      onClick: () => onNav('sites/new')
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 16
    }), " Create site")), /*#__PURE__*/React.createElement("div", {
      className: "card",
      style: {
        display: 'flex',
        alignItems: 'center',
        gap: 16,
        padding: 'var(--space-8) var(--space-10)'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "stat-icon",
      style: {
        background: 'var(--success-bg)',
        color: 'var(--success)'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "activity",
      size: 18
    })), /*#__PURE__*/React.createElement("div", {
      style: {
        flex: 1
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "ds-body",
      style: {
        fontWeight: 600
      }
    }, "All systems operational"), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption"
    }, "12 components healthy \xB7 auth, db, storage, proxy, deploy, functions\u2026")), /*#__PURE__*/React.createElement(Badge, {
      variant: "success",
      dot: true
    }, "Healthy")), /*#__PURE__*/React.createElement("div", {
      className: "stats-grid"
    }, stats.map(s => /*#__PURE__*/React.createElement("a", {
      key: s.label,
      className: "stat-card",
      href: "#",
      onClick: e => {
        e.preventDefault();
        onNav(s.to);
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "stat-card-top"
    }, /*#__PURE__*/React.createElement("div", {
      className: "stat-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: s.icon,
      size: 18
    })), /*#__PURE__*/React.createElement(Icon, {
      name: "chevron-right",
      size: 16,
      className: "dim"
    })), /*#__PURE__*/React.createElement("div", {
      className: "stat-count"
    }, s.count), /*#__PURE__*/React.createElement("div", {
      className: "stat-label"
    }, s.label)))), /*#__PURE__*/React.createElement("div", {
      className: "dash-cols",
      style: {
        display: 'grid',
        gridTemplateColumns: '1.4fr 1fr',
        gap: 20
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card"
    }, /*#__PURE__*/React.createElement("div", {
      className: "card-header"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Recent deployments"), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-link",
      onClick: () => onNav('sites')
    }, "View all")), /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        gap: 2
      }
    }, recent.map(d => /*#__PURE__*/React.createElement("div", {
      key: d.id,
      className: "row",
      style: {
        gap: 12,
        padding: '10px 0',
        borderBottom: '1px solid var(--border-default)'
      }
    }, /*#__PURE__*/React.createElement(StatusBadge, {
      status: d.status
    }), /*#__PURE__*/React.createElement("div", {
      style: {
        flex: 1,
        minWidth: 0
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "ds-body",
      style: {
        fontWeight: 550
      }
    }, d.site), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption",
      style: {
        overflow: 'hidden',
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap'
      }
    }, d.msg)), /*#__PURE__*/React.createElement("span", {
      className: "mono ds-caption"
    }, d.commit), /*#__PURE__*/React.createElement("span", {
      className: "ds-caption",
      style: {
        minWidth: 56,
        textAlign: 'right'
      }
    }, timeAgo(d.when)))))), /*#__PURE__*/React.createElement("div", {
      className: "card"
    }, /*#__PURE__*/React.createElement("div", {
      className: "card-header"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Jump back in")), /*#__PURE__*/React.createElement("div", {
      className: "col gap-2"
    }, [{
      l: 'Deploy a site from GitHub',
      i: 'rocket',
      to: 'sites/new'
    }, {
      l: 'Write a SQL query',
      i: 'terminal',
      to: 'sql'
    }, {
      l: 'Add a serverless function',
      i: 'box',
      to: 'functions'
    }, {
      l: 'Invite a teammate',
      i: 'users',
      to: 'users'
    }].map(q => /*#__PURE__*/React.createElement("a", {
      key: q.l,
      href: "#",
      onClick: e => {
        e.preventDefault();
        onNav(q.to);
      },
      className: "row",
      style: {
        gap: 12,
        padding: '11px 12px',
        borderRadius: 'var(--radius-s)',
        textDecoration: 'none',
        color: 'var(--fg-primary)',
        border: '1px solid var(--border-default)'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "stat-icon",
      style: {
        width: 30,
        height: 30
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: q.i,
      size: 15
    })), /*#__PURE__*/React.createElement("span", {
      className: "ds-body",
      style: {
        flex: 1,
        fontWeight: 500
      }
    }, q.l), /*#__PURE__*/React.createElement(Icon, {
      name: "chevron-right",
      size: 16,
      className: "dim"
    })))))));
  }

  /* ── Sites list ── */
  function SitesList({
    sites,
    onNav,
    loading
  }) {
    const [q, setQ] = useState('');
    const [filter, setFilter] = useState('all');
    const shown = sites.filter(s => (filter === 'all' || s.env === filter) && s.name.toLowerCase().includes(q.toLowerCase()));
    if (loading) {
      return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement(PageHeader, {
        title: "Sites"
      }), /*#__PURE__*/React.createElement("div", {
        className: "site-grid"
      }, [0, 1, 2].map(i => /*#__PURE__*/React.createElement(SkeletonCard, {
        key: i
      }))));
    }
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement(PageHeader, {
      title: "Sites",
      subtitle: "Deploy and host web apps straight from Git."
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      onClick: () => onNav('sites/new')
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 16
    }), " Create site")), sites.length === 0 ? /*#__PURE__*/React.createElement("div", {
      className: "empty-state"
    }, /*#__PURE__*/React.createElement("div", {
      className: "empty-state-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "rocket",
      size: 28
    })), /*#__PURE__*/React.createElement("div", {
      className: "empty-state-title"
    }, "Create your first site"), /*#__PURE__*/React.createElement("div", {
      className: "empty-state-text"
    }, "Connect a Git repository and BigBase will build, deploy, and serve it with a live preview URL \u2014 auto-redeploying on every push."), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      onClick: () => onNav('sites/new')
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 16
    }), " Create site")) : /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between',
        gap: 12
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 4,
        background: 'var(--bg-subtle)',
        padding: 3,
        borderRadius: 'var(--radius-s)'
      }
    }, ['all', 'production', 'preview'].map(f => /*#__PURE__*/React.createElement("button", {
      key: f,
      onClick: () => setFilter(f),
      className: "btn btn-sm",
      style: {
        background: filter === f ? 'var(--bg-surface)' : 'transparent',
        color: filter === f ? 'var(--fg-primary)' : 'var(--fg-secondary)',
        boxShadow: filter === f ? 'var(--shadow-xs)' : 'none',
        textTransform: 'capitalize',
        border: 'none'
      }
    }, f))), /*#__PURE__*/React.createElement("div", {
      className: "input-with-prefix",
      style: {
        width: 260
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "input-prefix"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "search",
      size: 14
    })), /*#__PURE__*/React.createElement("input", {
      className: "input",
      placeholder: "Search sites",
      value: q,
      onChange: e => setQ(e.target.value)
    }))), /*#__PURE__*/React.createElement("div", {
      className: "site-grid"
    }, shown.map(s => /*#__PURE__*/React.createElement("a", {
      key: s.id,
      className: "site-card",
      href: "#",
      onClick: e => {
        e.preventDefault();
        onNav('sites/' + s.id);
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "site-card-thumb",
      style: {
        background: THUMBS[s.thumb]
      }
    }, s.status === 'building' && /*#__PURE__*/React.createElement("span", {
      className: "badge badge-warning",
      style: {
        position: 'absolute',
        top: 10,
        right: 10
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "spinner spinner-sm",
      style: {
        width: 10,
        height: 10,
        borderWidth: 1.5
      }
    }), "Building"), s.status === 'failed' && /*#__PURE__*/React.createElement("span", {
      className: "badge badge-error",
      style: {
        position: 'absolute',
        top: 10,
        right: 10
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "dot"
    }), "Failed")), /*#__PURE__*/React.createElement("div", {
      className: "site-card-body"
    }, /*#__PURE__*/React.createElement("div", {
      className: "site-card-row"
    }, /*#__PURE__*/React.createElement("span", {
      className: "site-card-name"
    }, s.name), /*#__PURE__*/React.createElement(StatusBadge, {
      status: s.status
    })), /*#__PURE__*/React.createElement("span", {
      className: "site-card-url"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "globe",
      size: 12,
      style: {
        verticalAlign: '-2px',
        marginRight: 4
      }
    }), s.url), /*#__PURE__*/React.createElement("div", {
      className: "site-card-meta"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "git-branch",
      size: 13
    }), " ", s.branch, " \xB7 ", s.commit, " \xB7 ", timeAgo(s.updated))))))));
  }

  /* ── Site detail ── */
  function SiteDetail({
    site,
    onNav,
    onRedeploy
  }) {
    const [tab, setTab] = useState('overview');
    const tabs = ['overview', 'deployments', 'domains', 'logs', 'settings'];
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement("div", {
      className: "breadcrumb"
    }, /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: e => {
        e.preventDefault();
        onNav('sites');
      }
    }, "Sites"), /*#__PURE__*/React.createElement("span", {
      className: "sep"
    }, "/"), /*#__PURE__*/React.createElement("span", {
      style: {
        color: 'var(--fg-primary)'
      }
    }, site.name)), /*#__PURE__*/React.createElement("div", {
      className: "page-header",
      style: {
        marginBottom: 0
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 14
      }
    }, /*#__PURE__*/React.createElement("div", {
      style: {
        width: 44,
        height: 44,
        borderRadius: 'var(--radius-s)',
        background: THUMBS[site.thumb]
      }
    }), /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 10
      }
    }, /*#__PURE__*/React.createElement("h1", {
      className: "page-title"
    }, site.name), /*#__PURE__*/React.createElement(StatusBadge, {
      status: site.status
    })), /*#__PURE__*/React.createElement("a", {
      className: "site-card-url",
      href: "#",
      onClick: e => e.preventDefault(),
      style: {
        textDecoration: 'none'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "globe",
      size: 12,
      style: {
        verticalAlign: '-2px',
        marginRight: 4
      }
    }), site.url))), /*#__PURE__*/React.createElement("div", {
      className: "page-header-actions"
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary",
      onClick: () => onRedeploy(site)
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "refresh-cw",
      size: 15
    }), " Redeploy"), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "external-link",
      size: 15
    }), " Visit"))), /*#__PURE__*/React.createElement("div", {
      className: "tabs"
    }, tabs.map(t => /*#__PURE__*/React.createElement("button", {
      key: t,
      className: `tab ${tab === t ? 'active' : ''}`,
      style: {
        textTransform: 'capitalize'
      },
      onClick: () => setTab(t)
    }, t))), tab === 'overview' && /*#__PURE__*/React.createElement("div", {
      className: "dash-cols",
      style: {
        display: 'grid',
        gridTemplateColumns: '1.5fr 1fr',
        gap: 20
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card",
      style: {
        padding: 0,
        overflow: 'hidden'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "site-card-thumb",
      style: {
        height: 200,
        background: THUMBS[site.thumb]
      }
    }), /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        padding: 20
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Production deployment"), /*#__PURE__*/React.createElement(StatusBadge, {
      status: site.status
    })), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "globe",
      size: 15,
      className: "dim"
    }), /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: e => e.preventDefault(),
      className: "mono",
      style: {
        fontSize: 14,
        color: 'var(--fg-accent)',
        textDecoration: 'none'
      }
    }, site.url), /*#__PURE__*/React.createElement("button", {
      className: "toast-close"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "copy",
      size: 14
    }))), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption"
    }, "Deployed ", timeAgo(site.updated), " from ", /*#__PURE__*/React.createElement("span", {
      className: "mono"
    }, site.branch), " \xB7 ", site.framework))), /*#__PURE__*/React.createElement("div", {
      className: "card col gap-3"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Latest commit"), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "git-branch",
      size: 15,
      className: "dim"
    }), /*#__PURE__*/React.createElement("span", {
      className: "mono ds-body"
    }, site.commit)), /*#__PURE__*/React.createElement("div", {
      className: "ds-body",
      style: {
        fontWeight: 500
      }
    }, site.commitMsg || 'Initial deployment'), /*#__PURE__*/React.createElement("div", {
      className: "divider",
      style: {
        margin: '6px 0'
      }
    }), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between',
        gap: 12
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-caption",
      style: {
        flexShrink: 0
      }
    }, "Repository"), /*#__PURE__*/React.createElement("span", {
      className: "ds-body mono",
      style: {
        fontSize: 13,
        overflow: 'hidden',
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap'
      }
    }, site.repo)), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, "Branch"), /*#__PURE__*/React.createElement("span", {
      className: "ds-body mono",
      style: {
        fontSize: 13
      }
    }, site.branch)), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, "Root"), /*#__PURE__*/React.createElement("span", {
      className: "ds-body mono",
      style: {
        fontSize: 13
      }
    }, site.root)), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, "Environment"), /*#__PURE__*/React.createElement(Badge, {
      variant: "accent"
    }, site.env)))), tab === 'deployments' && /*#__PURE__*/React.createElement("div", {
      className: "table-wrap"
    }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Status"), /*#__PURE__*/React.createElement("th", null, "Commit"), /*#__PURE__*/React.createElement("th", null, "Branch"), /*#__PURE__*/React.createElement("th", null, "Duration"), /*#__PURE__*/React.createElement("th", null, "When"), /*#__PURE__*/React.createElement("th", null))), /*#__PURE__*/React.createElement("tbody", null, site.deployments.map(d => /*#__PURE__*/React.createElement("tr", {
      key: d.id
    }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(StatusBadge, {
      status: d.status
    })), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("span", {
      className: "mono"
    }, d.commit), " ", /*#__PURE__*/React.createElement("span", {
      className: "dim"
    }, d.msg)), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("span", {
      className: "mono"
    }, d.branch)), /*#__PURE__*/React.createElement("td", null, d.duration), /*#__PURE__*/React.createElement("td", null, timeAgo(d.when)), /*#__PURE__*/React.createElement("td", {
      className: "actions-cell"
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "more-horizontal",
      size: 16
    })))))))), tab === 'domains' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3"
    }, /*#__PURE__*/React.createElement("div", {
      className: "card row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 10
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "globe",
      size: 16,
      className: "dim"
    }), /*#__PURE__*/React.createElement("span", {
      className: "mono ds-body"
    }, site.url), /*#__PURE__*/React.createElement(Badge, {
      variant: "success",
      dot: true
    }, "Active")), /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, "BigBase subdomain \xB7 auto-issued TLS")), /*#__PURE__*/React.createElement("div", {
      className: "card col gap-3"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Add a custom domain"), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement("input", {
      className: "input",
      placeholder: "www.yourdomain.com"
    }), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary"
    }, "Add domain")), /*#__PURE__*/React.createElement("span", {
      className: "input-hint"
    }, "Point a CNAME at ", /*#__PURE__*/React.createElement("span", {
      className: "mono"
    }, site.url), " and BigBase issues a certificate automatically."))), tab === 'logs' && /*#__PURE__*/React.createElement("div", {
      className: "code-output",
      style: {
        maxHeight: 420
      }
    }, BUILD_LOG.map((l, i) => /*#__PURE__*/React.createElement("div", {
      className: "log-line",
      key: i
    }, /*#__PURE__*/React.createElement("span", {
      className: "log-time"
    }, l.t), /*#__PURE__*/React.createElement("span", {
      className: l.k === 'ok' ? 'log-ok' : '',
      style: {
        color: l.k === 'ok' ? undefined : '#d4d4d4'
      }
    }, l.m)))), tab === 'settings' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        maxWidth: 620
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card col gap-3"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Build settings"), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Framework preset"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      defaultValue: site.framework
    })), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Build command"), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      defaultValue: "npm run build"
    })), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Output directory"), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      defaultValue: "dist"
    })), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      style: {
        alignSelf: 'flex-start'
      }
    }, "Save changes")), /*#__PURE__*/React.createElement("div", {
      className: "card col gap-3",
      style: {
        borderColor: 'var(--error)'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title",
      style: {
        color: 'var(--error)'
      }
    }, "Danger zone"), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "ds-body",
      style: {
        fontWeight: 600
      }
    }, "Delete this site"), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption"
    }, "Permanently remove ", site.name, " and all deployments.")), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-danger"
    }, "Delete site")))));
  }

  /* ── Functions (list) ── */
  function Functions({
    onNav
  }) {
    const triggerIcon = {
      HTTP: 'zap',
      Schedule: 'clock',
      Event: 'activity'
    };
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement(PageHeader, {
      title: "Functions",
      subtitle: "Run server-side code on HTTP, schedule, or event triggers."
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      onClick: () => onNav('functions/' + FUNCTIONS[0].id)
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 16
    }), " Create function")), /*#__PURE__*/React.createElement("div", {
      className: "site-grid"
    }, FUNCTIONS.map(f => /*#__PURE__*/React.createElement("a", {
      key: f.id,
      className: "stat-card",
      href: "#",
      onClick: e => {
        e.preventDefault();
        onNav('functions/' + f.id);
      },
      style: {
        gap: 14
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "stat-card-top"
    }, /*#__PURE__*/React.createElement("div", {
      className: "stat-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "box",
      size: 18
    })), /*#__PURE__*/React.createElement(StatusBadge, {
      status: f.status === 'active' ? 'ready' : f.status
    })), /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        gap: 2
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "mono",
      style: {
        fontSize: 15,
        fontWeight: 600,
        color: 'var(--fg-primary)'
      }
    }, f.name), /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, f.runtime)), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 6,
        flexWrap: 'wrap'
      }
    }, /*#__PURE__*/React.createElement(Badge, {
      variant: "neutral"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: triggerIcon[f.trigger],
      size: 11,
      style: {
        verticalAlign: '-1px',
        marginRight: 2
      }
    }), f.trigger), /*#__PURE__*/React.createElement(Badge, {
      variant: "neutral"
    }, f.timeout, "s timeout")), /*#__PURE__*/React.createElement("div", {
      className: "divider",
      style: {
        margin: '2px 0'
      }
    }), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, "Created ", f.created), /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, "Updated ", timeAgo(f.updated)))))));
  }

  /* ── Data Studio (expanded: table browser + schema explorer + column ops) ── */
  function DataStudio() {
    const [active, setActive] = useState('users');
    const [view, setView] = useState('data');
    const schema = window.DATA.SCHEMA[active] || window.DATA.SCHEMA.users;
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement(PageHeader, {
      title: "Data Studio",
      subtitle: "Browse, edit, and shape your collections."
    }, view === 'data' ? /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 15
    }), " Add row") : /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 15
    }), " Add column")), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 20,
        alignItems: 'flex-start'
      }
    }, /*#__PURE__*/React.createElement("div", {
      style: {
        width: 200,
        flexShrink: 0
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "sidebar-section-title",
      style: {
        padding: '0 0 8px'
      }
    }, "Collections"), /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        gap: 1
      }
    }, COLLECTIONS.map(c => /*#__PURE__*/React.createElement("button", {
      key: c,
      onClick: () => setActive(c),
      className: "row",
      style: {
        gap: 8,
        padding: '8px 10px',
        borderRadius: 'var(--radius-s)',
        border: 'none',
        cursor: 'pointer',
        fontSize: 14,
        fontFamily: 'var(--font-mono)',
        background: active === c ? 'var(--brand-tint)' : 'transparent',
        color: active === c ? 'var(--fg-accent)' : 'var(--fg-secondary)'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "database",
      size: 14
    }), " ", c))), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm",
      style: {
        marginTop: 10,
        width: '100%',
        justifyContent: 'flex-start'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 13
    }), " New collection")), /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        flex: 1,
        gap: 14,
        minWidth: 0
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between',
        gap: 12
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 4,
        background: 'var(--bg-subtle)',
        padding: 3,
        borderRadius: 'var(--radius-s)'
      }
    }, [['data', 'Data'], ['schema', 'Schema']].map(([v, l]) => /*#__PURE__*/React.createElement("button", {
      key: v,
      onClick: () => setView(v),
      className: "btn btn-sm",
      style: {
        background: view === v ? 'var(--bg-surface)' : 'transparent',
        color: view === v ? 'var(--fg-primary)' : 'var(--fg-secondary)',
        boxShadow: view === v ? 'var(--shadow-xs)' : 'none',
        border: 'none'
      }
    }, l))), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-caption mono"
    }, active), /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, "\xB7 ", view === 'data' ? `${TABLE_ROWS.length} rows` : `${schema.length} columns`))), view === 'data' ? /*#__PURE__*/React.createElement("div", {
      className: "table-wrap"
    }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "id"), /*#__PURE__*/React.createElement("th", null, "email"), /*#__PURE__*/React.createElement("th", null, "role"), /*#__PURE__*/React.createElement("th", null, "verified"), /*#__PURE__*/React.createElement("th", null, "created"), /*#__PURE__*/React.createElement("th", null))), /*#__PURE__*/React.createElement("tbody", null, TABLE_ROWS.map(r => /*#__PURE__*/React.createElement("tr", {
      key: r.id
    }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("span", {
      className: "mono dim"
    }, r.id)), /*#__PURE__*/React.createElement("td", null, r.email), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(Badge, {
      variant: r.role === 'owner' ? 'accent' : 'neutral'
    }, r.role)), /*#__PURE__*/React.createElement("td", null, r.verified ? /*#__PURE__*/React.createElement(Badge, {
      variant: "success",
      dot: true
    }, "Yes") : /*#__PURE__*/React.createElement(Badge, {
      variant: "warning",
      dot: true
    }, "No")), /*#__PURE__*/React.createElement("td", {
      className: "dim"
    }, r.created), /*#__PURE__*/React.createElement("td", {
      className: "actions-cell"
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm",
      title: "Edit row"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "more-horizontal",
      size: 16
    })))))))) : /*#__PURE__*/React.createElement("div", {
      className: "table-wrap"
    }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Column"), /*#__PURE__*/React.createElement("th", null, "Type"), /*#__PURE__*/React.createElement("th", null, "Default"), /*#__PURE__*/React.createElement("th", null, "Nullable"), /*#__PURE__*/React.createElement("th", {
      style: {
        textAlign: 'right'
      }
    }, "Operations"))), /*#__PURE__*/React.createElement("tbody", null, schema.map(col => /*#__PURE__*/React.createElement("tr", {
      key: col.name
    }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "mono",
      style: {
        fontWeight: 600
      }
    }, col.name), col.pk && /*#__PURE__*/React.createElement(Badge, {
      variant: "accent"
    }, "PK"))), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("span", {
      className: "mono dim"
    }, col.type)), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("span", {
      className: "mono dim"
    }, col.def)), /*#__PURE__*/React.createElement("td", null, col.nullable ? /*#__PURE__*/React.createElement("span", {
      className: "dim"
    }, "nullable") : /*#__PURE__*/React.createElement("span", {
      className: "dim"
    }, "required")), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("div", {
      className: "actions-cell",
      style: {
        justifyContent: 'flex-end'
      }
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm"
    }, "Edit"), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm",
      style: {
        color: col.pk ? 'var(--fg-tertiary)' : 'var(--error)'
      },
      disabled: col.pk
    }, "Delete")))))))))));
  }

  /* ── Generic placeholder for un-built nav targets ── */
  function Placeholder({
    title
  }) {
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement(PageHeader, {
      title: title
    }), /*#__PURE__*/React.createElement("div", {
      className: "empty-state"
    }, /*#__PURE__*/React.createElement("div", {
      className: "empty-state-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "box",
      size: 28
    })), /*#__PURE__*/React.createElement("div", {
      className: "empty-state-title"
    }, title), /*#__PURE__*/React.createElement("div", {
      className: "empty-state-text"
    }, "This screen follows the same BigBase design system shown across Sites, Dashboard, and Data Studio. Wire it up next using the shared components.")));
  }
  Object.assign(window, {
    Login,
    Dashboard,
    SitesList,
    SiteDetail,
    Functions,
    DataStudio,
    Placeholder,
    PageHeader
  });
})();
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/admin-console/screens.jsx", error: String((e && e.message) || e) }); }

// ui_kits/admin-console/screens2.jsx
try { (() => {
/* BigBase Admin Console — secondary screens (Functions detail, Messaging, Settings) */
(function () {
  const {
    useState
  } = React;
  const Icon = window.Icon;
  const {
    FUNCTIONS,
    TEMPLATES,
    RUNTIMES,
    TRIGGERS,
    TEMPLATE_TYPES,
    timeAgo,
    statusVariant
  } = window.DATA;
  const {
    Badge,
    StatusBadge,
    Avatar,
    PageHeader
  } = window;

  /* ── Function detail ── */
  function FunctionDetail({
    fn,
    onNav
  }) {
    const [tab, setTab] = useState('code');
    const tabs = ['code', 'triggers', 'variables', 'logs'];
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement("div", {
      className: "breadcrumb"
    }, /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: e => {
        e.preventDefault();
        onNav('functions');
      }
    }, "Functions"), /*#__PURE__*/React.createElement("span", {
      className: "sep"
    }, "/"), /*#__PURE__*/React.createElement("span", {
      style: {
        color: 'var(--fg-primary)'
      },
      className: "mono"
    }, fn.name)), /*#__PURE__*/React.createElement("div", {
      className: "page-header",
      style: {
        marginBottom: 0
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 14
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "stat-icon",
      style: {
        width: 44,
        height: 44
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "box",
      size: 20
    })), /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 10
      }
    }, /*#__PURE__*/React.createElement("h1", {
      className: "page-title mono",
      style: {
        whiteSpace: 'nowrap'
      }
    }, fn.name), /*#__PURE__*/React.createElement(StatusBadge, {
      status: "ready"
    })), /*#__PURE__*/React.createElement("div", {
      className: "page-subtitle"
    }, fn.runtime, " \xB7 ", fn.trigger, " trigger \xB7 ", fn.invocations, " invocations \xB7 ", fn.errors, " errors"))), /*#__PURE__*/React.createElement("div", {
      className: "page-header-actions"
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "play",
      size: 15
    }), " Run now"), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "refresh-cw",
      size: 15
    }), " Deploy"))), /*#__PURE__*/React.createElement("div", {
      className: "tabs"
    }, tabs.map(t => /*#__PURE__*/React.createElement("button", {
      key: t,
      className: `tab ${tab === t ? 'active' : ''}`,
      style: {
        textTransform: 'capitalize'
      },
      onClick: () => setTab(t)
    }, t))), tab === 'code' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3"
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "index.js"), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-caption"
    }, "Deployed ", timeAgo(fn.deployed)), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "copy",
      size: 13
    }), " Copy"))), /*#__PURE__*/React.createElement("div", {
      className: "code-output",
      style: {
        maxHeight: 380
      }
    }, /*#__PURE__*/React.createElement("pre", {
      style: {
        color: '#d4d4d4'
      }
    }, fn.code))), tab === 'triggers' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        maxWidth: 620
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card col gap-4"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Trigger"), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Type"), /*#__PURE__*/React.createElement("select", {
      className: "input",
      defaultValue: fn.trigger
    }, TRIGGERS.map(t => /*#__PURE__*/React.createElement("option", {
      key: t
    }, t)))), fn.trigger === 'Schedule' ? /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Cron schedule"), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      defaultValue: fn.schedule || '0 6 * * *'
    }), /*#__PURE__*/React.createElement("span", {
      className: "input-hint"
    }, "Runs daily at 06:00 UTC. Uses standard cron syntax.")) : /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "HTTP endpoint"), /*#__PURE__*/React.createElement("div", {
      className: "input-with-prefix"
    }, /*#__PURE__*/React.createElement("span", {
      className: "input-prefix"
    }, "POST"), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      defaultValue: `https://fn.bigbase.local/${fn.name}`,
      readOnly: true
    })), /*#__PURE__*/React.createElement("span", {
      className: "input-hint"
    }, "Public endpoint. Add an API key in Variables to require auth.")), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Timeout (seconds)"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      type: "number",
      defaultValue: fn.timeout,
      style: {
        width: 120
      }
    })), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      style: {
        alignSelf: 'flex-start'
      }
    }, "Save changes"))), tab === 'variables' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        maxWidth: 620
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card col gap-4"
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Environment variables"), /*#__PURE__*/React.createElement(Badge, {
      variant: "neutral"
    }, fn.env.length)), fn.env.map((e, i) => /*#__PURE__*/React.createElement("div", {
      className: "row",
      key: i,
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      defaultValue: e.k
    }), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      defaultValue: e.v,
      type: e.v.includes('•') ? 'password' : 'text'
    }), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "x",
      size: 14
    })))), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary btn-sm",
      style: {
        alignSelf: 'flex-start'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 14
    }), " Add variable"))), tab === 'logs' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Execution history"), /*#__PURE__*/React.createElement("div", {
      className: "code-output",
      style: {
        maxHeight: 360
      }
    }, fn.logs.map((l, i) => /*#__PURE__*/React.createElement("div", {
      className: "log-line",
      key: i
    }, /*#__PURE__*/React.createElement("span", {
      className: "log-time"
    }, l.t), /*#__PURE__*/React.createElement("span", {
      className: l.k === 'ok' ? 'log-ok' : l.k === 'warn' ? 'log-warn' : '',
      style: {
        color: l.k ? undefined : '#d4d4d4'
      }
    }, l.m))))));
  }

  /* ── Messaging (list) ── */
  function Messaging({
    onNav
  }) {
    const [q, setQ] = useState('');
    const shown = TEMPLATES.filter(t => t.name.toLowerCase().includes(q.toLowerCase()));
    const typeIcon = {
      Email: 'mail',
      SMS: 'zap',
      Push: 'bell'
    };
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement(PageHeader, {
      title: "Messaging",
      subtitle: "Design templates for transactional email, SMS, and push."
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      onClick: () => onNav('messaging/' + TEMPLATES[0].id)
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 16
    }), " Create template")), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'flex-end'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "input-with-prefix",
      style: {
        width: 260
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "input-prefix"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "search",
      size: 14
    })), /*#__PURE__*/React.createElement("input", {
      className: "input",
      placeholder: "Search templates",
      value: q,
      onChange: e => setQ(e.target.value)
    }))), /*#__PURE__*/React.createElement("div", {
      className: "table-wrap"
    }, /*#__PURE__*/React.createElement("table", null, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", null, "Template"), /*#__PURE__*/React.createElement("th", null, "Type"), /*#__PURE__*/React.createElement("th", null, "Status"), /*#__PURE__*/React.createElement("th", null, "Sends"), /*#__PURE__*/React.createElement("th", null, "Updated"), /*#__PURE__*/React.createElement("th", null))), /*#__PURE__*/React.createElement("tbody", null, shown.map(t => /*#__PURE__*/React.createElement("tr", {
      key: t.id,
      style: {
        cursor: 'pointer'
      },
      onClick: () => onNav('messaging/' + t.id)
    }, /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: typeIcon[t.type],
      size: 15,
      className: "dim"
    }), /*#__PURE__*/React.createElement("span", {
      style: {
        fontWeight: 600
      }
    }, t.name))), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement(Badge, {
      variant: "neutral"
    }, t.type)), /*#__PURE__*/React.createElement("td", null, t.status === 'active' && /*#__PURE__*/React.createElement(Badge, {
      variant: "success",
      dot: true
    }, "Active"), t.status === 'draft' && /*#__PURE__*/React.createElement(Badge, {
      variant: "neutral",
      dot: true
    }, "Draft"), t.status === 'paused' && /*#__PURE__*/React.createElement(Badge, {
      variant: "warning",
      dot: true
    }, "Paused")), /*#__PURE__*/React.createElement("td", {
      className: "dim"
    }, t.sends), /*#__PURE__*/React.createElement("td", {
      className: "dim"
    }, timeAgo(t.updated)), /*#__PURE__*/React.createElement("td", {
      className: "actions-cell"
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "chevron-right",
      size: 16
    })))))))));
  }

  /* ── Messaging detail (template editor + preview) ── */
  function MessagingDetail({
    tpl,
    onNav
  }) {
    const [subject, setSubject] = useState(tpl.subject);
    const [body, setBody] = useState(tpl.body);
    const sample = {
      name: 'Daniel',
      workspace: 'Acme',
      site: 'marketing-site',
      duration: '38s',
      url: 'marketing-site.bigbase.local',
      code: '418-203',
      reset_url: 'bigbase.local/r/abc',
      inviter: 'Maya',
      accept_url: 'bigbase.local/i/xyz'
    };
    const fill = s => s.replace(/\{\{(\w+)\}\}/g, (_, k) => sample[k] || `{{${k}}}`);
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement("div", {
      className: "breadcrumb"
    }, /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: e => {
        e.preventDefault();
        onNav('messaging');
      }
    }, "Messaging"), /*#__PURE__*/React.createElement("span", {
      className: "sep"
    }, "/"), /*#__PURE__*/React.createElement("span", {
      style: {
        color: 'var(--fg-primary)'
      }
    }, tpl.name)), /*#__PURE__*/React.createElement("div", {
      className: "page-header",
      style: {
        marginBottom: 0
      }
    }, /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 10
      }
    }, /*#__PURE__*/React.createElement("h1", {
      className: "page-title"
    }, tpl.name), /*#__PURE__*/React.createElement(Badge, {
      variant: "neutral"
    }, tpl.type)), /*#__PURE__*/React.createElement("div", {
      className: "page-subtitle"
    }, tpl.sends, " sends \xB7 updated ", timeAgo(tpl.updated))), /*#__PURE__*/React.createElement("div", {
      className: "page-header-actions"
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary"
    }, "Send test"), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary"
    }, "Save changes"))), /*#__PURE__*/React.createElement("div", {
      className: "dash-cols",
      style: {
        display: 'grid',
        gridTemplateColumns: '1fr 1fr',
        gap: 20,
        alignItems: 'flex-start'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card col gap-4"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Editor"), tpl.type === 'Email' && /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Subject"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      value: subject,
      onChange: e => setSubject(e.target.value)
    })), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Body"), /*#__PURE__*/React.createElement("textarea", {
      className: "input",
      style: {
        minHeight: 200,
        lineHeight: 1.6
      },
      value: body,
      onChange: e => setBody(e.target.value)
    })), /*#__PURE__*/React.createElement("div", {
      className: "col gap-2"
    }, /*#__PURE__*/React.createElement("span", {
      className: "input-label"
    }, "Variables"), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 6,
        flexWrap: 'wrap'
      }
    }, tpl.vars.map(v => /*#__PURE__*/React.createElement("span", {
      key: v,
      className: "badge badge-accent mono",
      style: {
        textTransform: 'none'
      }
    }, `{{${v}}}`))), /*#__PURE__*/React.createElement("span", {
      className: "input-hint"
    }, "Click a variable to insert it. Values are filled per recipient at send time."))), /*#__PURE__*/React.createElement("div", {
      className: "card col gap-3",
      style: {
        background: 'var(--bg-subtle)'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Preview"), /*#__PURE__*/React.createElement("div", {
      className: "card",
      style: {
        background: 'var(--bg-surface)',
        boxShadow: 'var(--shadow-xs)'
      }
    }, tpl.type === 'Email' && /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("div", {
      className: "ds-caption"
    }, "Subject"), /*#__PURE__*/React.createElement("div", {
      className: "ds-body",
      style: {
        fontWeight: 600,
        marginBottom: 14
      }
    }, fill(subject)), /*#__PURE__*/React.createElement("div", {
      className: "divider",
      style: {
        margin: '0 0 14px'
      }
    })), /*#__PURE__*/React.createElement("div", {
      className: "ds-body",
      style: {
        whiteSpace: 'pre-wrap',
        lineHeight: 1.7
      }
    }, fill(body))), /*#__PURE__*/React.createElement("span", {
      className: "input-hint"
    }, "Rendered with sample values for ", /*#__PURE__*/React.createElement("span", {
      className: "mono"
    }, tpl.vars.join(', ')), "."))));
  }

  /* ── Settings (Account / Workspace / Billing) ── */
  function Settings({
    user
  }) {
    const [tab, setTab] = useState('account');
    const tabs = [['account', 'Account'], ['workspace', 'Workspace'], ['billing', 'Billing']];
    const members = [{
      email: user.email,
      role: 'Owner',
      initial: user.email[0].toUpperCase()
    }, {
      email: 'maya@acme.io',
      role: 'Admin',
      initial: 'M'
    }, {
      email: 'sam@studio.co',
      role: 'Member',
      initial: 'S'
    }];
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement(PageHeader, {
      title: "Settings",
      subtitle: "Manage your account, workspace, and plan."
    }), /*#__PURE__*/React.createElement("div", {
      className: "tabs"
    }, tabs.map(([id, label]) => /*#__PURE__*/React.createElement("button", {
      key: id,
      className: `tab ${tab === id ? 'active' : ''}`,
      onClick: () => setTab(id)
    }, label))), tab === 'account' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        maxWidth: 620
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card col gap-4"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Profile"), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Email"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      defaultValue: user.email
    })), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Display name"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      defaultValue: user.email.split('@')[0]
    })), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      style: {
        alignSelf: 'flex-start'
      }
    }, "Save changes")), /*#__PURE__*/React.createElement("div", {
      className: "card col gap-4"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Password"), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Current password"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      type: "password",
      defaultValue: "demo1234"
    })), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "New password"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      type: "password",
      placeholder: "At least 6 characters"
    })), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary",
      style: {
        alignSelf: 'flex-start'
      }
    }, "Update password")), /*#__PURE__*/React.createElement("div", {
      className: "card row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "ds-body",
      style: {
        fontWeight: 600
      }
    }, "Two-factor authentication"), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption"
    }, "Add a second step at sign-in with an authenticator app.")), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "check-circle",
      size: 15
    }), " Enable 2FA"))), tab === 'workspace' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        maxWidth: 720
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card col gap-4"
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Workspace"), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Workspace name"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      defaultValue: "Acme"
    })), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Workspace URL"), /*#__PURE__*/React.createElement("div", {
      className: "input-with-prefix"
    }, /*#__PURE__*/React.createElement("span", {
      className: "input-prefix"
    }, "bigbase.local/"), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      defaultValue: "acme"
    }))), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      style: {
        alignSelf: 'flex-start'
      }
    }, "Save changes")), /*#__PURE__*/React.createElement("div", {
      className: "card col gap-3"
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "card-title"
    }, "Members"), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary btn-sm"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 13
    }), " Invite member")), /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        gap: 1
      }
    }, members.map(m => /*#__PURE__*/React.createElement("div", {
      key: m.email,
      className: "row",
      style: {
        gap: 12,
        padding: '10px 0',
        borderBottom: '1px solid var(--border-default)'
      }
    }, /*#__PURE__*/React.createElement(Avatar, {
      email: m.email,
      size: 32
    }), /*#__PURE__*/React.createElement("span", {
      className: "ds-body",
      style: {
        flex: 1
      }
    }, m.email), /*#__PURE__*/React.createElement(Badge, {
      variant: m.role === 'Owner' ? 'accent' : 'neutral'
    }, m.role)))))), tab === 'billing' && /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        maxWidth: 720
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "card row",
      style: {
        justifyContent: 'space-between',
        alignItems: 'center'
      }
    }, /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 10
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-h3"
    }, "Self-hosted"), /*#__PURE__*/React.createElement(Badge, {
      variant: "success",
      dot: true
    }, "Active")), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption",
      style: {
        marginTop: 4
      }
    }, "Single-binary instance \xB7 unlimited projects \xB7 no usage caps")), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary",
      disabled: true
    }, "Manage plan")), /*#__PURE__*/React.createElement("div", {
      className: "empty-state",
      style: {
        padding: 'var(--space-20) var(--space-16)'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "empty-state-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "info",
      size: 28
    })), /*#__PURE__*/React.createElement("div", {
      className: "empty-state-title"
    }, "No invoices yet"), /*#__PURE__*/React.createElement("div", {
      className: "empty-state-text"
    }, "You're running a self-hosted BigBase instance, so there's nothing to bill. Hosted plans with usage-based billing are coming soon."))));
  }
  Object.assign(window, {
    FunctionDetail,
    Messaging,
    MessagingDetail,
    Settings
  });
})();
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/admin-console/screens2.jsx", error: String((e && e.message) || e) }); }

// ui_kits/admin-console/ui.jsx
try { (() => {
/* BigBase Admin Console — shared UI primitives + sidebar */
(function () {
  const {
    useState,
    useEffect,
    createContext,
    useContext,
    useCallback
  } = React;
  const Icon = window.Icon;

  /* ── Badge ── */
  function Badge({
    variant = 'neutral',
    dot,
    children
  }) {
    return /*#__PURE__*/React.createElement("span", {
      className: `badge badge-${variant}`
    }, dot && /*#__PURE__*/React.createElement("span", {
      className: "dot"
    }), children);
  }
  function StatusBadge({
    status
  }) {
    const v = window.DATA.statusVariant(status);
    const label = status.charAt(0).toUpperCase() + status.slice(1);
    return /*#__PURE__*/React.createElement(Badge, {
      variant: v,
      dot: status !== 'ready'
    }, status === 'building' ? /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("span", {
      className: "spinner spinner-sm",
      style: {
        width: 10,
        height: 10,
        borderWidth: 1.5
      }
    }), label) : /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("span", {
      className: "dot"
    }), label));
  }

  /* ── Avatar ── */
  function Avatar({
    email,
    size = 30
  }) {
    const initial = (email || '?')[0].toUpperCase();
    return /*#__PURE__*/React.createElement("div", {
      className: "avatar",
      style: {
        width: size,
        height: size,
        fontSize: size * 0.42
      }
    }, initial);
  }

  /* ── Preview banner ── */
  function PreviewBanner() {
    return /*#__PURE__*/React.createElement("div", {
      className: "preview-banner"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "alert-triangle",
      size: 16
    }), "Preview mode \u2014 showing mock data. Connect a backend to deploy for real.");
  }

  /* ── Toasts ── */
  const ToastCtx = createContext(null);
  function useToast() {
    return useContext(ToastCtx);
  }
  function ToastProvider({
    children
  }) {
    const [toasts, setToasts] = useState([]);
    const push = useCallback(t => {
      const id = Math.random().toString(36).slice(2);
      setToasts(p => [...p, {
        id,
        ...t
      }]);
      setTimeout(() => setToasts(p => p.filter(x => x.id !== id)), t.duration || 4200);
    }, []);
    const remove = id => setToasts(p => p.filter(x => x.id !== id));
    const iconFor = {
      success: 'check-circle',
      error: 'alert-triangle',
      info: 'info'
    };
    return /*#__PURE__*/React.createElement(ToastCtx.Provider, {
      value: push
    }, children, /*#__PURE__*/React.createElement("div", {
      className: "toast-container"
    }, toasts.map(t => /*#__PURE__*/React.createElement("div", {
      key: t.id,
      className: `toast toast-${t.type || 'info'}`
    }, /*#__PURE__*/React.createElement(Icon, {
      name: iconFor[t.type || 'info'],
      size: 18,
      className: "toast-icon"
    }), /*#__PURE__*/React.createElement("div", {
      className: "toast-body"
    }, /*#__PURE__*/React.createElement("div", {
      className: "toast-title"
    }, t.title), t.msg && /*#__PURE__*/React.createElement("div", {
      className: "toast-msg"
    }, t.msg)), /*#__PURE__*/React.createElement("button", {
      className: "toast-close",
      onClick: () => remove(t.id)
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "x",
      size: 15
    }))))));
  }

  /* ── Skeleton ── */
  function SkeletonCard() {
    return /*#__PURE__*/React.createElement("div", {
      className: "card",
      style: {
        display: 'flex',
        flexDirection: 'column',
        gap: 12
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "skeleton skeleton-card",
      style: {
        height: 110
      }
    }), /*#__PURE__*/React.createElement("div", {
      className: "skeleton skeleton-text",
      style: {
        width: '70%'
      }
    }), /*#__PURE__*/React.createElement("div", {
      className: "skeleton skeleton-text",
      style: {
        width: '45%'
      }
    }));
  }

  /* ── Theme selector ── */
  function ThemeSelect({
    theme,
    onTheme
  }) {
    const [open, setOpen] = useState(false);
    const ref = React.useRef(null);
    const THEMES = window.DATA.THEMES;
    const current = THEMES.find(t => t.key === theme) || THEMES[0];
    useEffect(() => {
      if (!open) return;
      const close = e => {
        if (ref.current && !ref.current.contains(e.target)) setOpen(false);
      };
      const esc = e => {
        if (e.key === 'Escape') setOpen(false);
      };
      document.addEventListener('mousedown', close);
      document.addEventListener('keydown', esc);
      return () => {
        document.removeEventListener('mousedown', close);
        document.removeEventListener('keydown', esc);
      };
    }, [open]);
    return /*#__PURE__*/React.createElement("div", {
      className: "theme-select-wrap",
      ref: ref
    }, /*#__PURE__*/React.createElement("button", {
      className: "theme-select",
      "aria-haspopup": "listbox",
      "aria-expanded": open,
      onClick: () => setOpen(o => !o)
    }, /*#__PURE__*/React.createElement("span", {
      className: "theme-swatch",
      style: {
        background: current.swatch
      }
    }), /*#__PURE__*/React.createElement("span", {
      style: {
        flex: 1
      }
    }, current.month === 'Default' ? 'Indigo' : `${current.month} · ${current.label}`), /*#__PURE__*/React.createElement(Icon, {
      name: "chevron-down",
      size: 15,
      className: "dim"
    })), open && /*#__PURE__*/React.createElement("div", {
      className: "theme-menu",
      role: "listbox",
      "aria-label": "Accent theme"
    }, THEMES.map(t => /*#__PURE__*/React.createElement("button", {
      key: t.key,
      role: "option",
      "aria-selected": t.key === theme,
      className: `theme-menu-item ${t.key === theme ? 'selected' : ''}`,
      onClick: () => {
        onTheme(t.key);
        setOpen(false);
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "theme-dot",
      style: {
        background: t.swatch
      }
    }), /*#__PURE__*/React.createElement("span", {
      style: {
        flex: 1
      }
    }, t.month === 'Default' ? 'Indigo' : t.month), t.key === theme && /*#__PURE__*/React.createElement(Icon, {
      name: "check",
      size: 13
    })))));
  }

  /* ── Sidebar ── */
  const NAV = [{
    section: 'Overview',
    items: [{
      id: 'dashboard',
      label: 'Dashboard',
      icon: 'layout-dashboard'
    }]
  }, {
    section: 'Build',
    items: [{
      id: 'sites',
      label: 'Sites',
      icon: 'rocket'
    }, {
      id: 'functions',
      label: 'Functions',
      icon: 'box'
    }]
  }, {
    section: 'Data',
    items: [{
      id: 'data',
      label: 'Data Studio',
      icon: 'database'
    }, {
      id: 'sql',
      label: 'SQL Editor',
      icon: 'terminal'
    }, {
      id: 'storage',
      label: 'Storage',
      icon: 'hard-drive'
    }]
  }, {
    section: 'Auth',
    items: [{
      id: 'users',
      label: 'Users',
      icon: 'users'
    }]
  }, {
    section: 'Engage',
    items: [{
      id: 'messaging',
      label: 'Messaging',
      icon: 'mail'
    }]
  }, {
    section: 'DevOps',
    items: [{
      id: 'repos',
      label: 'Git Repos',
      icon: 'git-branch'
    }, {
      id: 'cici',
      label: 'CI / CD',
      icon: 'git-pull-request'
    }, {
      id: 'monitoring',
      label: 'Monitoring',
      icon: 'activity'
    }]
  }];
  function Sidebar({
    route,
    onNav,
    user,
    dark,
    onToggleDark,
    onLogout,
    theme,
    onTheme
  }) {
    const active = id => route === id || route.startsWith(id + '/');
    return /*#__PURE__*/React.createElement("nav", {
      className: "sidebar"
    }, /*#__PURE__*/React.createElement("div", {
      className: "sidebar-logo"
    }, /*#__PURE__*/React.createElement("img", {
      className: "sidebar-logo-mark",
      src: "../../assets/bigbase-logo.svg",
      alt: ""
    }), /*#__PURE__*/React.createElement("span", {
      className: "sidebar-logo-text"
    }, "BigBase")), NAV.map(grp => /*#__PURE__*/React.createElement("div", {
      className: "sidebar-section",
      key: grp.section
    }, /*#__PURE__*/React.createElement("div", {
      className: "sidebar-section-title"
    }, grp.section), /*#__PURE__*/React.createElement("ul", {
      className: "sidebar-nav"
    }, grp.items.map(it => /*#__PURE__*/React.createElement("li", {
      key: it.id
    }, /*#__PURE__*/React.createElement("a", {
      href: "#",
      className: active(it.id) ? 'active' : '',
      onClick: e => {
        e.preventDefault();
        onNav(it.id);
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "sidebar-nav-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: it.icon,
      size: 18
    })), /*#__PURE__*/React.createElement("span", null, it.label))))))), /*#__PURE__*/React.createElement("div", {
      className: "sidebar-spacer"
    }), /*#__PURE__*/React.createElement("div", {
      className: "sidebar-section"
    }, /*#__PURE__*/React.createElement("div", {
      className: "sidebar-section-title"
    }, "Appearance"), /*#__PURE__*/React.createElement(ThemeSelect, {
      theme: theme,
      onTheme: onTheme
    }), /*#__PURE__*/React.createElement("ul", {
      className: "sidebar-nav",
      style: {
        marginTop: 'var(--space-3)'
      }
    }, /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: e => {
        e.preventDefault();
        onToggleDark();
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "sidebar-nav-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: dark ? 'sun' : 'moon',
      size: 18
    })), /*#__PURE__*/React.createElement("span", null, dark ? 'Light mode' : 'Dark mode'))), /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement("a", {
      href: "#",
      className: route === 'settings' ? 'active' : '',
      onClick: e => {
        e.preventDefault();
        onNav('settings');
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "sidebar-nav-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "settings",
      size: 18
    })), /*#__PURE__*/React.createElement("span", null, "Settings"))))), /*#__PURE__*/React.createElement("div", {
      className: "sidebar-footer"
    }, /*#__PURE__*/React.createElement("div", {
      className: "sidebar-user"
    }, /*#__PURE__*/React.createElement(Avatar, {
      email: user.email
    }), /*#__PURE__*/React.createElement("div", {
      style: {
        minWidth: 0,
        flex: 1
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "sidebar-email",
      style: {
        fontWeight: 600,
        color: 'var(--fg-primary)'
      }
    }, user.email.split('@')[0]), /*#__PURE__*/React.createElement("div", {
      className: "sidebar-email"
    }, user.email)), /*#__PURE__*/React.createElement("button", {
      className: "toast-close",
      title: "Log out",
      onClick: onLogout
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "log-out",
      size: 16
    })))));
  }
  Object.assign(window, {
    Badge,
    StatusBadge,
    Avatar,
    PreviewBanner,
    ToastProvider,
    useToast,
    SkeletonCard,
    Sidebar,
    ThemeSelect
  });
})();
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/admin-console/ui.jsx", error: String((e && e.message) || e) }); }

// ui_kits/admin-console/wizard.jsx
try { (() => {
/* BigBase Admin Console — Create Site wizard (Journey A) */
(function () {
  const {
    useState,
    useEffect,
    useRef
  } = React;
  const Icon = window.Icon;
  const {
    REPOS,
    FRAMEWORKS,
    BUILD_LOG,
    THUMBS
  } = window.DATA;
  function WizardRail({
    step
  }) {
    const steps = ['Source', 'Configure', 'Deploy'];
    return /*#__PURE__*/React.createElement("div", {
      className: "wizard-steps"
    }, steps.map((s, i) => /*#__PURE__*/React.createElement(React.Fragment, {
      key: s
    }, /*#__PURE__*/React.createElement("div", {
      className: `wizard-step ${i < step ? 'done' : ''} ${i === step ? 'active' : ''}`
    }, /*#__PURE__*/React.createElement("div", {
      className: "wizard-step-num"
    }, i < step ? /*#__PURE__*/React.createElement(Icon, {
      name: "check",
      size: 14
    }) : i + 1), /*#__PURE__*/React.createElement("div", {
      className: "wizard-step-label"
    }, s)), i < steps.length - 1 && /*#__PURE__*/React.createElement("div", {
      className: `wizard-step-line ${i < step ? 'done' : ''}`
    }))));
  }

  /* ── Step 1: Source ── */
  function SourceStep({
    source,
    setSource,
    gitConnected,
    setGitConnected,
    selectedRepo,
    setSelectedRepo,
    search,
    setSearch
  }) {
    const filtered = REPOS.filter(r => r.name.toLowerCase().includes(search.toLowerCase()));
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "ds-h3",
      style: {
        marginBottom: 6
      }
    }, "Where's your code?"), /*#__PURE__*/React.createElement("div", {
      className: "ds-body dim"
    }, "Choose a source. BigBase detects the framework and builds it for you.")), /*#__PURE__*/React.createElement("div", {
      className: "choice-grid"
    }, /*#__PURE__*/React.createElement("div", {
      className: `choice-card ${source === 'github' ? 'selected' : ''}`,
      onClick: () => setSource('github')
    }, /*#__PURE__*/React.createElement("div", {
      className: "choice-card-check"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "check",
      size: 18
    })), /*#__PURE__*/React.createElement("div", {
      className: "choice-card-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "github",
      size: 22
    })), /*#__PURE__*/React.createElement("div", {
      className: "choice-card-title"
    }, "Connect Git"), /*#__PURE__*/React.createElement("div", {
      className: "choice-card-desc"
    }, "Deploy from a GitHub repository. Auto-redeploys on every push.")), /*#__PURE__*/React.createElement("div", {
      className: `choice-card ${source === 'bigbase' ? 'selected' : ''}`,
      onClick: () => setSource('bigbase')
    }, /*#__PURE__*/React.createElement("div", {
      className: "choice-card-check"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "check",
      size: 18
    })), /*#__PURE__*/React.createElement("div", {
      className: "choice-card-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "box",
      size: 22
    })), /*#__PURE__*/React.createElement("div", {
      className: "choice-card-title"
    }, "Existing BigBase repo"), /*#__PURE__*/React.createElement("div", {
      className: "choice-card-desc"
    }, "Use a repo already hosted on this BigBase instance.")), /*#__PURE__*/React.createElement("div", {
      className: "choice-card disabled"
    }, /*#__PURE__*/React.createElement("div", {
      className: "choice-card-icon"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "rocket",
      size: 22
    })), /*#__PURE__*/React.createElement("div", {
      className: "choice-card-title"
    }, "Start from template"), /*#__PURE__*/React.createElement("div", {
      className: "choice-card-desc"
    }, "Astro, Next.js & more."), /*#__PURE__*/React.createElement("span", {
      className: "badge badge-neutral",
      style: {
        alignSelf: 'flex-start',
        marginTop: 4
      }
    }, "Coming soon"))), source === 'github' && /*#__PURE__*/React.createElement("div", {
      className: "card",
      style: {
        marginTop: 4
      }
    }, !gitConnected ? /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        alignItems: 'center',
        gap: 14,
        padding: '20px 0'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "choice-card-icon",
      style: {
        width: 52,
        height: 52
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "github",
      size: 26
    })), /*#__PURE__*/React.createElement("div", {
      className: "ds-h3"
    }, "Connect your GitHub account"), /*#__PURE__*/React.createElement("div", {
      className: "ds-body dim",
      style: {
        textAlign: 'center',
        maxWidth: 360
      }
    }, "Authorize BigBase to read repositories and receive push events."), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      onClick: () => setGitConnected(true)
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "github",
      size: 16
    }), " Authorize GitHub")) : /*#__PURE__*/React.createElement("div", {
      className: "col gap-3"
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "github",
      size: 16
    }), /*#__PURE__*/React.createElement("span", {
      className: "ds-body",
      style: {
        fontWeight: 600
      }
    }, "danielvm"), /*#__PURE__*/React.createElement("span", {
      className: "badge badge-success"
    }, /*#__PURE__*/React.createElement("span", {
      className: "dot"
    }), "Connected")), /*#__PURE__*/React.createElement("div", {
      className: "input-with-prefix",
      style: {
        width: 240
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "input-prefix"
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "search",
      size: 14
    })), /*#__PURE__*/React.createElement("input", {
      className: "input",
      placeholder: "Search repos",
      value: search,
      onChange: e => setSearch(e.target.value)
    }))), /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        gap: 1,
        border: '1px solid var(--border-default)',
        borderRadius: 'var(--radius-s)',
        overflow: 'hidden'
      }
    }, filtered.map(r => /*#__PURE__*/React.createElement("div", {
      key: r.name,
      onClick: () => setSelectedRepo(r.name),
      style: {
        display: 'flex',
        alignItems: 'center',
        gap: 12,
        padding: '12px 14px',
        cursor: 'pointer',
        background: selectedRepo === r.name ? 'var(--brand-tint)' : 'var(--bg-surface)'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "git-branch",
      size: 16,
      className: "dim"
    }), /*#__PURE__*/React.createElement("div", {
      style: {
        flex: 1,
        minWidth: 0
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        gap: 6
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-body",
      style: {
        fontWeight: 600
      }
    }, r.name), r.private && /*#__PURE__*/React.createElement("span", {
      className: "badge badge-neutral"
    }, "Private")), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption"
    }, r.desc, " \xB7 ", r.lang, " \xB7 ", r.updated)), selectedRepo === r.name ? /*#__PURE__*/React.createElement("span", {
      className: "badge badge-accent"
    }, "Selected") : /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary btn-sm"
    }, "Select"))), filtered.length === 0 && /*#__PURE__*/React.createElement("div", {
      className: "dim",
      style: {
        padding: 16,
        fontSize: 14
      }
    }, "No repos match \"", search, "\".")))), source === 'bigbase' && /*#__PURE__*/React.createElement("div", {
      className: "card",
      style: {
        marginTop: 4
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "col",
      style: {
        gap: 1,
        border: '1px solid var(--border-default)',
        borderRadius: 'var(--radius-s)',
        overflow: 'hidden'
      }
    }, REPOS.slice(0, 3).map(r => /*#__PURE__*/React.createElement("div", {
      key: r.name,
      onClick: () => setSelectedRepo(r.name),
      style: {
        display: 'flex',
        alignItems: 'center',
        gap: 12,
        padding: '12px 14px',
        cursor: 'pointer',
        background: selectedRepo === r.name ? 'var(--brand-tint)' : 'var(--bg-surface)'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "box",
      size: 16,
      className: "dim"
    }), /*#__PURE__*/React.createElement("div", {
      style: {
        flex: 1
      }
    }, /*#__PURE__*/React.createElement("span", {
      className: "ds-body",
      style: {
        fontWeight: 600
      }
    }, r.name), /*#__PURE__*/React.createElement("div", {
      className: "ds-caption"
    }, "on this instance")), selectedRepo === r.name ? /*#__PURE__*/React.createElement("span", {
      className: "badge badge-accent"
    }, "Selected") : /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary btn-sm"
    }, "Select"))))));
  }

  /* ── Step 2: Configure ── */
  function ConfigureStep({
    cfg,
    setCfg
  }) {
    const set = (k, v) => setCfg(p => ({
      ...p,
      [k]: v
    }));
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "ds-h3",
      style: {
        marginBottom: 6
      }
    }, "Configure your site"), /*#__PURE__*/React.createElement("div", {
      className: "ds-body dim"
    }, "We detected ", /*#__PURE__*/React.createElement("b", {
      style: {
        color: 'var(--fg-primary)'
      }
    }, cfg.framework), ". Adjust if needed.")), /*#__PURE__*/React.createElement("div", {
      className: "card col gap-4"
    }, /*#__PURE__*/React.createElement("div", {
      className: "form-row",
      style: {
        display: 'flex',
        gap: 16
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "input-group",
      style: {
        flex: 1
      }
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Site name"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      value: cfg.name,
      onChange: e => set('name', e.target.value.replace(/[^a-z0-9-]/g, ''))
    }), /*#__PURE__*/React.createElement("span", {
      className: "input-hint"
    }, "Used for your URL: ", /*#__PURE__*/React.createElement("span", {
      className: "mono"
    }, cfg.name || 'your-site', ".bigbase.local"))), /*#__PURE__*/React.createElement("div", {
      className: "input-group",
      style: {
        width: 200
      }
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Production branch"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      value: cfg.branch,
      onChange: e => set('branch', e.target.value)
    }))), /*#__PURE__*/React.createElement("div", {
      className: "form-row",
      style: {
        display: 'flex',
        gap: 16
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "input-group",
      style: {
        flex: 1
      }
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Framework preset"), /*#__PURE__*/React.createElement("select", {
      className: "input",
      value: cfg.framework,
      onChange: e => {
        const f = FRAMEWORKS.find(x => x.name === e.target.value);
        setCfg(p => ({
          ...p,
          framework: f.name,
          build: f.build,
          install: f.install,
          output: f.output
        }));
      }
    }, FRAMEWORKS.map(f => /*#__PURE__*/React.createElement("option", {
      key: f.id,
      value: f.name
    }, f.name)))), /*#__PURE__*/React.createElement("div", {
      className: "input-group",
      style: {
        width: 200
      }
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Root directory"), /*#__PURE__*/React.createElement("input", {
      className: "input",
      value: cfg.root,
      onChange: e => set('root', e.target.value)
    })))), /*#__PURE__*/React.createElement("details", {
      className: "card",
      style: {
        padding: 0
      }
    }, /*#__PURE__*/React.createElement("summary", {
      style: {
        padding: 'var(--space-8) var(--space-10)',
        cursor: 'pointer',
        fontWeight: 600,
        fontSize: 14,
        display: 'flex',
        alignItems: 'center',
        gap: 8,
        listStyle: 'none'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "chevron-right",
      size: 16,
      className: "dim"
    }), " Build & output settings"), /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        padding: '0 var(--space-10) var(--space-10)'
      }
    }, /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Install command"), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      value: cfg.install,
      onChange: e => set('install', e.target.value)
    })), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Build command"), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      value: cfg.build,
      onChange: e => set('build', e.target.value)
    })), /*#__PURE__*/React.createElement("div", {
      className: "input-group"
    }, /*#__PURE__*/React.createElement("label", {
      className: "input-label"
    }, "Output directory"), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      value: cfg.output,
      onChange: e => set('output', e.target.value)
    })))), /*#__PURE__*/React.createElement("details", {
      className: "card",
      style: {
        padding: 0
      }
    }, /*#__PURE__*/React.createElement("summary", {
      style: {
        padding: 'var(--space-8) var(--space-10)',
        cursor: 'pointer',
        fontWeight: 600,
        fontSize: 14,
        display: 'flex',
        alignItems: 'center',
        gap: 8,
        listStyle: 'none'
      }
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "chevron-right",
      size: 16,
      className: "dim"
    }), " Environment variables ", /*#__PURE__*/React.createElement("span", {
      className: "badge badge-neutral",
      style: {
        marginLeft: 4
      }
    }, cfg.env.length)), /*#__PURE__*/React.createElement("div", {
      className: "col gap-3",
      style: {
        padding: '0 var(--space-10) var(--space-10)'
      }
    }, cfg.env.map((e, i) => /*#__PURE__*/React.createElement("div", {
      className: "row",
      key: i,
      style: {
        gap: 8
      }
    }, /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      placeholder: "KEY",
      value: e.k,
      onChange: ev => setCfg(p => {
        const env = [...p.env];
        env[i] = {
          ...env[i],
          k: ev.target.value
        };
        return {
          ...p,
          env
        };
      })
    }), /*#__PURE__*/React.createElement("input", {
      className: "input mono",
      placeholder: "value",
      value: e.v,
      onChange: ev => setCfg(p => {
        const env = [...p.env];
        env[i] = {
          ...env[i],
          v: ev.target.value
        };
        return {
          ...p,
          env
        };
      })
    }), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost btn-sm",
      onClick: () => setCfg(p => ({
        ...p,
        env: p.env.filter((_, j) => j !== i)
      }))
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "x",
      size: 14
    })))), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary btn-sm",
      style: {
        alignSelf: 'flex-start'
      },
      onClick: () => setCfg(p => ({
        ...p,
        env: [...p.env, {
          k: '',
          v: ''
        }]
      }))
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "plus",
      size: 14
    }), " Add variable"))));
  }

  /* ── Step 3: Deploy ── */
  function DeployStep({
    cfg,
    onReady
  }) {
    const [lines, setLines] = useState([]);
    const [done, setDone] = useState(false);
    const logRef = useRef(null);
    useEffect(() => {
      let i = 0;
      const t = setInterval(() => {
        if (i >= BUILD_LOG.length) {
          clearInterval(t);
          setDone(true);
          onReady && onReady();
          return;
        }
        setLines(prev => [...prev, BUILD_LOG[i]]);
        i++;
      }, 520);
      return () => clearInterval(t);
    }, []);
    useEffect(() => {
      if (logRef.current) logRef.current.scrollTop = logRef.current.scrollHeight;
    }, [lines]);
    const pct = Math.min(100, Math.round(lines.length / BUILD_LOG.length * 100));
    return /*#__PURE__*/React.createElement("div", {
      className: "col gap-4"
    }, /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between'
      }
    }, /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "ds-h3",
      style: {
        marginBottom: 6
      }
    }, done ? 'Your site is live 🎉' : 'Building your site…'), /*#__PURE__*/React.createElement("div", {
      className: "ds-body dim"
    }, cfg.name, ".bigbase.local \xB7 ", cfg.framework, " \xB7 ", cfg.branch)), done ? /*#__PURE__*/React.createElement("span", {
      className: "badge badge-success"
    }, /*#__PURE__*/React.createElement("span", {
      className: "dot"
    }), "Ready") : /*#__PURE__*/React.createElement("span", {
      className: "badge badge-warning"
    }, /*#__PURE__*/React.createElement("span", {
      className: "spinner spinner-sm",
      style: {
        width: 10,
        height: 10,
        borderWidth: 1.5
      }
    }), "Building")), /*#__PURE__*/React.createElement("div", {
      className: "bar-track",
      style: {
        height: 6,
        background: 'var(--bg-subtle)',
        borderRadius: 999,
        overflow: 'hidden'
      }
    }, /*#__PURE__*/React.createElement("div", {
      style: {
        height: '100%',
        width: pct + '%',
        background: done ? 'var(--success)' : 'var(--brand-500)',
        borderRadius: 999,
        transition: 'width .4s var(--ease-emphasized)'
      }
    })), /*#__PURE__*/React.createElement("div", {
      className: "code-output",
      ref: logRef,
      style: {
        maxHeight: 280
      }
    }, lines.map((l, i) => /*#__PURE__*/React.createElement("div", {
      className: "log-line",
      key: i
    }, /*#__PURE__*/React.createElement("span", {
      className: "log-time"
    }, l.t), /*#__PURE__*/React.createElement("span", {
      className: l.k === 'ok' ? 'log-ok' : '',
      style: {
        color: l.k === 'ok' ? undefined : '#d4d4d4'
      }
    }, l.m))), !done && /*#__PURE__*/React.createElement("div", {
      className: "log-line"
    }, /*#__PURE__*/React.createElement("span", {
      className: "log-time"
    }, "\xA0"), /*#__PURE__*/React.createElement("span", {
      style: {
        color: 'var(--brand-100)'
      }
    }, "\u258D"))));
  }

  /* ── Wizard shell ── */
  function CreateSite({
    onCancel,
    onDeployed,
    onBackToList
  }) {
    const [step, setStep] = useState(0);
    const [source, setSource] = useState(null);
    const [gitConnected, setGitConnected] = useState(false);
    const [selectedRepo, setSelectedRepo] = useState(null);
    const [search, setSearch] = useState('');
    const [deployReady, setDeployReady] = useState(false);
    const [cfg, setCfg] = useState({
      name: '',
      branch: 'main',
      framework: 'Astro',
      root: '/',
      install: 'npm install',
      build: 'npm run build',
      output: 'dist',
      env: []
    });
    useEffect(() => {
      if (selectedRepo && !cfg.name) {
        const guessed = selectedRepo.split('/')[1] || 'my-site';
        setCfg(p => ({
          ...p,
          name: guessed
        }));
      }
    }, [selectedRepo]);
    const canNext = step === 0 ? !!selectedRepo : step === 1 ? cfg.name.length > 1 : false;
    const next = () => setStep(s => Math.min(2, s + 1));
    const back = () => setStep(s => Math.max(0, s - 1));
    const finish = () => {
      const site = {
        id: 'ste_' + cfg.name,
        name: cfg.name,
        framework: cfg.framework,
        repo: selectedRepo,
        branch: cfg.branch,
        root: cfg.root,
        status: 'ready',
        url: cfg.name + '.bigbase.local',
        commit: 'a3f91c2',
        commitMsg: 'Initial deployment',
        updated: Date.now(),
        env: 'production',
        thumb: 'grad-emerald',
        deployments: [{
          id: 'dep_new',
          status: 'ready',
          branch: cfg.branch,
          commit: 'a3f91c2',
          msg: 'Initial deployment',
          when: Date.now(),
          duration: '38s'
        }]
      };
      onDeployed(site);
    };
    return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("div", {
      className: "breadcrumb"
    }, /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: e => {
        e.preventDefault();
        onBackToList();
      }
    }, "Sites"), /*#__PURE__*/React.createElement("span", {
      className: "sep"
    }, "/"), /*#__PURE__*/React.createElement("span", null, "Create site")), /*#__PURE__*/React.createElement("div", {
      className: "page-header"
    }, /*#__PURE__*/React.createElement("h1", {
      className: "page-title"
    }, "Create a new site"), /*#__PURE__*/React.createElement("button", {
      className: "btn btn-ghost",
      onClick: onCancel
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "x",
      size: 16
    }), " Cancel")), /*#__PURE__*/React.createElement("div", {
      style: {
        maxWidth: 760
      }
    }, /*#__PURE__*/React.createElement(WizardRail, {
      step: step
    }), step === 0 && /*#__PURE__*/React.createElement(SourceStep, {
      source,
      setSource,
      gitConnected,
      setGitConnected,
      selectedRepo,
      setSelectedRepo,
      search,
      setSearch
    }), step === 1 && /*#__PURE__*/React.createElement(ConfigureStep, {
      cfg: cfg,
      setCfg: setCfg
    }), step === 2 && /*#__PURE__*/React.createElement(DeployStep, {
      cfg: cfg,
      onReady: () => setDeployReady(true)
    }), /*#__PURE__*/React.createElement("div", {
      className: "row",
      style: {
        justifyContent: 'space-between',
        marginTop: 28
      }
    }, /*#__PURE__*/React.createElement("button", {
      className: "btn btn-secondary",
      onClick: step === 0 ? onCancel : back,
      disabled: step === 2
    }, /*#__PURE__*/React.createElement(Icon, {
      name: "arrow-left",
      size: 16
    }), " ", step === 0 ? 'Cancel' : 'Back'), step < 2 && /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      disabled: !canNext,
      onClick: next
    }, step === 1 ? /*#__PURE__*/React.createElement(React.Fragment, null, "Deploy ", /*#__PURE__*/React.createElement(Icon, {
      name: "rocket",
      size: 16
    })) : /*#__PURE__*/React.createElement(React.Fragment, null, "Continue ", /*#__PURE__*/React.createElement(Icon, {
      name: "arrow-right",
      size: 16
    }))), step === 2 && /*#__PURE__*/React.createElement("button", {
      className: "btn btn-primary",
      disabled: !deployReady,
      onClick: finish
    }, deployReady ? /*#__PURE__*/React.createElement(React.Fragment, null, "View site ", /*#__PURE__*/React.createElement(Icon, {
      name: "arrow-right",
      size: 16
    })) : /*#__PURE__*/React.createElement(React.Fragment, null, /*#__PURE__*/React.createElement("span", {
      className: "spinner spinner-sm"
    }), " Building\u2026")))));
  }
  window.CreateSite = CreateSite;
})();
})(); } catch (e) { __ds_ns.__errors.push({ path: "ui_kits/admin-console/wizard.jsx", error: String((e && e.message) || e) }); }

if (__ds_scope.__ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1$13xveek === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1$13xveek = __ds_scope.__ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1 !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_LoginPage_1dmvnk1 : __ds_scope.LoginPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_DashboardPage_120zej4$mogxob === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_DashboardPage_120zej4$mogxob = __ds_scope.__ds_default_source_bigbase_ui_pages_DashboardPage_120zej4 !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_DashboardPage_120zej4 : __ds_scope.DashboardPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6$txrdad === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6$txrdad = __ds_scope.__ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6 !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_DataStudioPage_136wfp6 : __ds_scope.DataStudioPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr$1reqlu2 === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr$1reqlu2 = __ds_scope.__ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_SqlEditorPage_7q50pr : __ds_scope.SqlEditorPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui$1n11yp1 === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui$1n11yp1 = __ds_scope.__ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_UsersPage_1wq27ui : __ds_scope.UsersPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5$15vgx4w === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5$15vgx4w = __ds_scope.__ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5 !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_GitReposPage_vkrrw5 : __ds_scope.GitReposPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_DeployPage_17gbuph$7njolc === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_DeployPage_17gbuph$7njolc = __ds_scope.__ds_default_source_bigbase_ui_pages_DeployPage_17gbuph !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_DeployPage_17gbuph : __ds_scope.DeployPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza$1579x4h === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza$1579x4h = __ds_scope.__ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_MessagingPage_1kjsdza : __ds_scope.MessagingPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p$mt7puw === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p$mt7puw = __ds_scope.__ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_StoragePage_1vp8g1p : __ds_scope.StoragePage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx$1socyh8 === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx$1socyh8 = __ds_scope.__ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_FunctionsPage_8zrdcx : __ds_scope.FunctionsPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf$8ghtly === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf$8ghtly = __ds_scope.__ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_ForgePage_i5i2rf : __ds_scope.ForgePage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c$n2l73f === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c$n2l73f = __ds_scope.__ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_CiciPage_1lsz65c : __ds_scope.CiciPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse$1r6cjcp === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse$1r6cjcp = __ds_scope.__ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_MonitoringPage_1edjse : __ds_scope.MonitoringPage;
if (__ds_scope.__ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t$m6vebk === undefined) __ds_scope.__ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t$m6vebk = __ds_scope.__ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_pages_NotFoundPage_bw692t : __ds_scope.NotFoundPage;
if (__ds_scope.__ds_default_source_bigbase_ui_Layout_2g8xmy$1x3lec === undefined) __ds_scope.__ds_default_source_bigbase_ui_Layout_2g8xmy$1x3lec = __ds_scope.__ds_default_source_bigbase_ui_Layout_2g8xmy !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_Layout_2g8xmy : __ds_scope.Layout;
if (__ds_scope.__ds_default_source_bigbase_ui_App_tku48d$1urd5a3 === undefined) __ds_scope.__ds_default_source_bigbase_ui_App_tku48d$1urd5a3 = __ds_scope.__ds_default_source_bigbase_ui_App_tku48d !== undefined ? __ds_scope.__ds_default_source_bigbase_ui_App_tku48d : __ds_scope.App;

__ds_ns.App = __ds_scope.App;

__ds_ns.Layout = __ds_scope.Layout;

__ds_ns.Avatar = __ds_scope.Avatar;

__ds_ns.Badge = __ds_scope.Badge;

__ds_ns.StatusBadge = __ds_scope.StatusBadge;

__ds_ns.Button = __ds_scope.Button;

__ds_ns.Card = __ds_scope.Card;

__ds_ns.CardHeader = __ds_scope.CardHeader;

__ds_ns.EmptyState = __ds_scope.EmptyState;

__ds_ns.Input = __ds_scope.Input;

__ds_ns.PageHeader = __ds_scope.PageHeader;

__ds_ns.Spinner = __ds_scope.Spinner;

__ds_ns.Tabs = __ds_scope.Tabs;

__ds_ns.ThemeSelect = __ds_scope.ThemeSelect;

__ds_ns.ThemeProvider = __ds_scope.ThemeProvider;

__ds_ns.THEMES = __ds_scope.THEMES;

__ds_ns.THEME_BY_KEY = __ds_scope.THEME_BY_KEY;

__ds_ns.CiciPage = __ds_scope.CiciPage;

__ds_ns.DashboardPage = __ds_scope.DashboardPage;

__ds_ns.DataStudioPage = __ds_scope.DataStudioPage;

__ds_ns.DeployPage = __ds_scope.DeployPage;

__ds_ns.ForgePage = __ds_scope.ForgePage;

__ds_ns.FunctionsPage = __ds_scope.FunctionsPage;

__ds_ns.GitReposPage = __ds_scope.GitReposPage;

__ds_ns.LoginPage = __ds_scope.LoginPage;

__ds_ns.MessagingPage = __ds_scope.MessagingPage;

__ds_ns.MonitoringPage = __ds_scope.MonitoringPage;

__ds_ns.NotFoundPage = __ds_scope.NotFoundPage;

__ds_ns.SqlEditorPage = __ds_scope.SqlEditorPage;

__ds_ns.StoragePage = __ds_scope.StoragePage;

__ds_ns.UsersPage = __ds_scope.UsersPage;

})();
