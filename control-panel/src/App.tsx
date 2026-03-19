import "./App.css";
import AppHeader from "./components/AppHeader";
import AppRoutes from "./Routes";

function App() {
  return (
    <div className="app-root">
      <AppHeader />
      <main className="app-content">
        <AppRoutes />
      </main>
    </div>
  );
}

export default App;
