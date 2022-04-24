class CustomElement extends HTMLElement {
    constructor() {
        super();
        if (this.templateContent) {
            const shadowRoot = this.attachShadow({ mode: 'open' });
            shadowRoot.appendChild(this.templateContent.cloneNode(true));
        }
    }

    clearChildren(element) {
        while(element.firstChild) {
            element.removeChild(element.firstChild);
        }
    }
}

export default CustomElement