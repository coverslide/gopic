'use strict';
(function() {

	function updateHistory(url) {
		window.history.pushState('', '', url);
		window.dispatchEvent(new Event("locationchange", { url }));
	}

	function defineTemplate(html){
		const template = document.createElement('template');
		template.innerHTML = html;
		return template;
	}

	class CustomElement extends HTMLElement {
		static template = defineTemplate('');

		getTemplate(){
			return CustomElement.template;
		}

		constructor() {
			super();
			const templateContent = this.getTemplate().content;

			const shadowRoot = this.attachShadow({mode: 'open'});
			shadowRoot.appendChild(templateContent.cloneNode(true));
		}

		clearChildren(root) {
			while(root.firstChild) {
				root.removeChild(root.firstChild);
			}
		}
	}

	class MainApp extends CustomElement {

		static template = defineTemplate('<gopic-app></gopic-app>');

		getTemplate(){
			return MainApp.template;
		}

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
			this.shadowRoot.querySelector('gopic-app').setAttribute("path", path);
		}
	}

	class GopicApp extends CustomElement {

		static template = defineTemplate(`
<bread-crumb></bread-crumb>
<file-list></file-list>
`);

		getTemplate(){
			return GopicApp.template;
		}

		static get observedAttributes() { return ['path']; }

		attributeChangedCallback(name, oldValue, newValue) {
			this.shadowRoot.querySelector('bread-crumb').setAttribute("path", newValue);
			this.shadowRoot.querySelector('file-list').setAttribute("path", newValue);
		}
	}

	class BreadCrumb extends CustomElement {

		static get observedAttributes() { return ['path']; }

		static template = defineTemplate(`
<nav></nav>
`);

		getTemplate(){
			return BreadCrumb.template;
		}

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
				console.log(pathPart, currentPath, currentPart);

				if (currentPart === 0) {
					const link = document.createElement('a');
					link.setAttribute("href", "/");
					link.appendChild(document.createTextNode("Home"));
					nav.appendChild(link);
				} else if (pathParts.length > currentPart) {
					const separator = document.createElement('span');
					separator.appendChild(document.createTextNode("/"));
					nav.appendChild(separator);
					if (pathParts.length > currentPart+1) {
						const link = document.createElement('a');
						currentPath += "/" + pathPart;
						link.setAttribute("href", currentPath);
						link.addEventListener('click', (event) => {
							event.preventDefault();
							updateHistory(event.target.getAttribute('href'));
						})
						link.appendChild(document.createTextNode(decodeURI(pathPart)));
						nav.appendChild(link);
					} else {
						nav.appendChild(document.createTextNode(decodeURI(pathPart)));
					}
				}
				currentPart += 1;
			}
		}
	}

	class FileList extends CustomElement {
		static get observedAttributes() { return ['path']; }

		static template = defineTemplate(`
<section></section>
`);

		getTemplate(){
			return FileList.template;
		}

		attributeChangedCallback(name, oldValue, newValue) {
			if (name == "path") {
				this.onPathChanged(oldValue,  newValue);
			}
		}

		onPathChanged(oldValue,  newValue){
			fetch(newValue + '?json=true').then(response => response.json())
			.then(data => this.populate(newValue, data))
			.catch(err => this.handleError(err));
		}

		handleError(err){
			console.error(err);
			this.clearChildren();
			this.shadowRoot.appendChild(document.createTextNode(err.message))
		}

		clearChildren() {
			super.clearChildren(this.shadowRoot.querySelector('section'));
		}

		populate(path, folderData){
			this.clearChildren();

			const files = folderData
			.filter(file => !file.filename.match(/^\./))
			.sort((a, b) => {
				if (a.isDir && !b.isDir) {
						return -1;
				} else if (b.isDir && !a.isDir) {
					return 1;
				}
				return a.filename > b.filename ? 1 : -1;
			});

			const root = this.shadowRoot.querySelector('section');
			for (let folder of files) {
				const entry = document.createElement("file-entry");
				entry.setAttribute("path", path);
				entry.setAttribute("filename", folder.filename);
				entry.setAttribute("directory", folder.isDir);
				root.appendChild(entry);
			}
		}
	}

	class FileEntry extends CustomElement {
		static template = defineTemplate(`
<style>
img {
	height: 50px;
	width: 50px;
}
</style>
`);

		getTemplate(){
			return FileEntry.template;
		}

		connectedCallback(){
			const path = this.getAttribute("path").replace(/^\/$/, '');
			const filename = this.getAttribute("filename");
			const directory = this.getAttribute("directory") === "true";
			const body = document.createElement('div');
			const link = document.createElement('a');
			link.setAttribute('href', `${path}/${filename}`);
			if (directory) {
				link.appendChild(document.createElement('folder-icon'));
				link.addEventListener('click', (event) => {
					event.preventDefault();
					updateHistory(event.target.getAttribute('href'));
				})
			} else {
				const img = document.createElement('img');
				img.setAttribute("src", `${path}/${filename}?thumbnail=true`);
				link.appendChild(img);
			}
			link.appendChild(document.createTextNode(filename));
			body.appendChild(link);
			this.shadowRoot.appendChild(body);
		}
	}

	class FolderIcon extends CustomElement {
		static template = defineTemplate(`
<svg width="50" height="50" viewBox="0 0 490 490">
	<polygon fill="wheat" points="410.3,447.2 0,447.2 79.7,157.9 490,157.9 			"/>
	<polygon fill="tan" points="62.2,134.9 410.3,134.9 410.3,90.6 205.3,90.6 184.7,42.8 0,42.8 0,360.9"/>
</svg>
`);

		getTemplate(){
			return FolderIcon.template;
		}
	}

	customElements.define("main-app", MainApp);
	customElements.define("gopic-app", GopicApp);
	customElements.define("bread-crumb", BreadCrumb);
	customElements.define("file-list", FileList);
	customElements.define("file-entry", FileEntry);
	customElements.define("folder-icon", FolderIcon);

	const mainApp = document.createElement("main-app");
	document.body.appendChild(mainApp);
})();
