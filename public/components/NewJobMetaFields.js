class NewJobMetaFields extends HTMLElement {
  connectedCallback() {
    this.appendChild(this.createNewMetaInputs());
    this.appendChild(this.createNewButton());
  }

  disconnectedCallback() {
    const newButton = this.querySelector("#add-meta");
    newButton.removeEventListener("click", (e) => this.onAddButtonClicked(e));

    const removeButtons = this.querySelectorAll("button[id$='-remove']");
    removeButtons.forEach((button) => {
      button.removeEventListener("click", (e) => this.onRemoveButtonClicked(e));
    });
  }

  createNewMetaInputs(meta) {
    if (!meta) {
      meta = {
        key: "",
        value: "",
      };
    }

    const metaDiv = document.createElement("div");
    metaDiv.id = `${meta.key}-meta`;
    const idx = this.querySelectorAll("div[id$='-meta']").length;
    metaDiv.dataset.metaIdx = idx;
    {
      const keyField = document.createElement("input");
      keyField.type = "text";
      keyField.value = meta.key ?? "";
      keyField.placeholder = "Key";
      keyField.name = `metas[${idx}].key`;
      metaDiv.appendChild(keyField);

      const valueField = document.createElement("input");
      valueField.type = "text";
      valueField.value = meta.value ?? "";
      valueField.placeholder = "Value";
      valueField.name = `metas[${idx}].value`;
      metaDiv.appendChild(valueField);
    }

    const removeButton = document.createElement("button");
    removeButton.id = `${meta.key}-remove`;
    removeButton.type = "button";
    removeButton.dataset.metaKey = meta.key;
    removeButton.innerText = "-";
    removeButton.addEventListener("click", (e) => {
      this.onRemoveButtonClicked(e);
    });
    metaDiv.appendChild(removeButton);

    return metaDiv;
  }

  createNewButton() {
    const button = document.createElement("button");
    button.type = "button";
    button.id = "add-meta";
    button.innerText = "+";
    button.addEventListener("click", (e) => {
      this.onAddButtonClicked(e);
    });
    return button;
  }

  /**
   * @param {Mouseevent} event
   */
  onAddButtonClicked(event) {
    const metaDiv = this.createNewMetaInputs();

    const addMetaButton = this.querySelector("#add-meta");
    addMetaButton.insertAdjacentElement("beforebegin", metaDiv);
  }

  /**
   * @param {MouseEvent} event
   * @returns {void}
   */
  onRemoveButtonClicked(event) {
    event.target.removeEventListener("click", (e) =>
      this.onRemoveButtonClicked(e)
    );
    const metaKey = event.target.dataset.metaKey;
    const metaDiv = document.getElementById(`${metaKey}-meta`);

    // get all next %-meta divs
    let nextMetaDiv = metaDiv.nextElementSibling;
    while (nextMetaDiv && nextMetaDiv.tagName === "DIV") {
      const idx = Number(nextMetaDiv.dataset.metaIdx);
      nextMetaDiv.dataset.metaIdx = idx - 1;
      nextMetaDiv.querySelector("input[name$='.key']").name = `metas[${
        idx - 1
      }].key`;
      nextMetaDiv.querySelector("input[name$='.value']").name = `metas[${
        idx - 1
      }].value`;
      nextMetaDiv = nextMetaDiv.nextElementSibling;
    }

    metaDiv.remove();
  }
}

customElements.define("new-job-meta-fields", NewJobMetaFields);
