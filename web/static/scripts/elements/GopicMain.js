import CustomElement from "./CustomElement.js";

const TEMPLATE = `
<style>
    section {
        width: 1024px;
        margin-left: auto;
        margin-right: auto;
    }
</style>
<section>
    <bread-crumb></bread-crumb>
    <file-list></file-list>
</section>
`;

class GopicMain extends CustomElement {
    static tag = "gopic-main";
    static template = TEMPLATE;

    static get observedAttributes() { return ['path']; }

    attributeChangedCallback(name, oldValue, newValue) {
        this.shadowRoot.querySelector('bread-crumb').setAttribute("path", newValue);
        this.shadowRoot.querySelector('file-list').setAttribute("path", newValue);
    }
}

export default GopicMain;