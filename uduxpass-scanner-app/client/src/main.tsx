import { createRoot } from "react-dom/client";
import App from "./App";
import "./index.css";
import { registerServiceWorker } from "./lib/registerServiceWorker";

// Render app
createRoot(document.getElementById("root")!).render(<App />);

// Register service worker for PWA functionality
if (import.meta.env.PROD) {
  registerServiceWorker()
    .then((result) => {
      if (result.success) {
        console.log('[PWA] App is ready for offline use');
      }
    })
    .catch((error) => {
      console.error('[PWA] Service worker registration error:', error);
    });
}
