import CustomElement from "./CustomElement.js";

class MainApp extends CustomElement {
    static tag = "main-app";
    static template = "<gopic-main></gopic-main>";

    connectedCallback() {
        const path = window.location.pathname;
        this.onPathUpdated(path);
        window.addEventListener('popstate', (event) => {
            const path = window.location.pathname;
            this.onPathUpdated(path);
        });
        window.addEventListener('locationchange', (event) => {
            const path = window.location.pathname;
            this.onPathUpdated(path);
        });
    }

    onPathUpdated(path) {
        this.shadowRoot.querySelector('gopic-main').setAttribute("path", path);
    }
}

export default MainApp;