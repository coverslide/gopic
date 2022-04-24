import CustomElement from "./CustomElement.js";

const TEMPLATE = `<style>
section.loading {
	opacity: .5;
}
</style>
<section></section>
`;

class FileList extends CustomElement {
	static tag = "file-list";
	static template = TEMPLATE;

	static get observedAttributes() { return ['path']; }

	attributeChangedCallback(name, oldValue, newValue) {
		if (name == "path") {
			this.onPathChanged(oldValue, newValue);
		}
	}

	startLoading() {
		this.shadowRoot.querySelector('section').classList.add("loading");
	}

	stopLoading() {
		this.shadowRoot.querySelector('section').classList.remove("loading");
	}

	onPathChanged(oldValue, newValue) {
		this.startLoading();
		fetch(newValue + '?json=true').then(response => response.json())
			.then(data => this.populate(newValue, data))
			.catch(err => this.handleError(err))
			.then(() => this.stopLoading());
	}

	handleError(err) {
		console.error(err);
		this.clearChildren();
		this.shadowRoot.appendChild(document.createTextNode(err.message))
	}

	clearChildren() {
		super.clearChildren(this.shadowRoot.querySelector('section'));
	}

	populate(path, folderData) {
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

export default FileList;