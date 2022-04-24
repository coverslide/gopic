import updateHistory from "../helpers/updateHistory.js";
import CustomElement from "./CustomElement.js";

const TEMPLATE = `<style>
a, a:link {
    font-weight: bold;
    font-family: sans-serif;
    color: #2c8898;
    text-decoration: none;
}
a:hover {
    text-decoration: underline;
}
</style>
<a></a>`;

class BreadCrumbLink extends CustomElement {
    static tag = "breadcrumb-link";
    static template = TEMPLATE;
    static get observedAttributes() { return ['path', 'name']; }

    connectedCallback() {
        const link = this.shadowRoot.querySelector("a");
        link.href = this.getAttribute("href");
        link.appendChild(document.createTextNode(this.getAttribute("name")));
        link.addEventListener('click', (event) => {
            event.preventDefault();
            updateHistory(event.target.getAttribute('href'));
        })
    }
}

export default BreadCrumbLink;