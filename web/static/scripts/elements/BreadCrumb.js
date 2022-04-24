import updateHistory from "../helpers/updateHistory.js";
import CustomElement from "./CustomElement.js";

const TEMPLATE = `<nav></nav>`;

class BreadCrumb extends CustomElement {
	static tag = "bread-crumb";
	static template = TEMPLATE;

	static get observedAttributes() { return ['path']; }

	attributeChangedCallback(name, oldValue, newValue) {
		if (name == 'path') {
			this.onPathChanged(oldValue, newValue)
		}
	}

	clearChildren() {
		super.clearChildren(this.shadowRoot.querySelector('nav'));
	}

	onPathChanged(oldValue, newValue) {
		this.clearChildren();
		const nav = this.shadowRoot.querySelector('nav');
		const pathParts = newValue.split('/');
		let currentPath = '';
		let currentPart = 0;
		for (let pathPart of pathParts) {
			if (currentPart === 0) {
				const link = document.createElement('breadcrumb-link');
				link.setAttribute("href", "/");
				link.setAttribute("name", "Home");
				nav.appendChild(link);
			} else if (pathParts.length > currentPart) {
				const separator = document.createElement('span');
				separator.appendChild(document.createTextNode("/"));
				nav.appendChild(separator);
				if (pathParts.length > currentPart + 1) {
					const link = document.createElement('breadcrumb-link');
					currentPath += "/" + pathPart;
					link.setAttribute("href", currentPath);
					link.setAttribute("name", decodeURI(pathPart));
					nav.appendChild(link);
				} else {
					nav.appendChild(document.createTextNode(decodeURI(pathPart)));
				}
			}
			currentPart += 1;
		}
	}
}

export default BreadCrumb;