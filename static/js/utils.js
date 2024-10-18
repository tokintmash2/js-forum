export const $ = (q) => {
  const elements = document.querySelectorAll(q);
  if (elements.length > 1) {
    return elements;
  } else if (elements.length === 1) {
    return elements[0];
  }
};

export const validateEmail = (email) => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(String(email).toLowerCase());
};

// Returns a DOM element, spcified by inputs as follows
// type: div, button
// name: id
// classNames: class names
// contents: innerHTML
// link: link
export function makeElement(
  type,
  name,
  classNames,
  contents,
  link,
  attributes,
  children = []
) {
  const element = document.createElement(type);
  if (name) {
    element.id = name;
  }
  if (classNames) {
    if (Array.isArray(classNames)) {
      element.classList.add(...classNames);
    } else {
      element.classList.add(classNames);
    }
  }
  if (contents) {
    element.innerHTML = contents;
  }
  if (link) {
    element.href = link;
  }

  if (attributes) {
    for (const [key, value] of Object.entries(attributes)) {
      element.setAttribute(key, value);
    }
  }

  if (children) {
    children.forEach((child) => element.appendChild(child));
  }

  return element;
}

/**
 * @typedef {Object} ElementData
 * @property {string} type - The type of the element (e.g., 'div', 'a').
 * @property {string} [name] - The id of the element.
 * @property {string|string[]} [classNames] - The class or classes of the element.
 * @property {string} [contents] - The innerHTML of the element.
 * @property {string} [link] - The href of the element if it is a link.
 */

/**
 *
 * @param {ElementData} data
 * @returns {HTMLElement}
 */

export function makeElements(data) {
  const element = document.createElement(data.type);
  if (data.name) {
    element.id = data.name;
  }
  if (data.classNames) {
    if (Array.isArray(data.classNames)) {
      element.classList.add(...data.classNames);
    } else {
      element.classList.add(data.classNames);
    }
  }
  if (data.contents) {
    element.innerHTML = data.contents;
  }
  if (data.link) {
    element.href = data.link;
  }

  if (data.attributes) {
    for (const [key, value] of Object.entries(data.attributes)) {
      element.setAttribute(key, value);
    }
  }

  if (data.children) {
    for (const child of data.children) {
      element.appendChild(child);
    }
  }

  return element;
}

export function clearPage() {
  $("#app").innerHTML = "";
}

export function clearMainContentContainer() {
  let mainContent = $(".main-content");
  if (mainContent) {
    mainContent.innerHTML = "";
  } else {
    mainContent = makeElements({ type: "div", classNames: "main-content" });
    $("#app").appendChild(mainContent);
  }
}

export function formatTimestamp(timestamp) {
  const date = new Date(timestamp);

  // Format the date to dd.mm.yyyy
  const formattedDate = date.toLocaleDateString("et-EE");

  // Format the time to hh:mm
  const formattedTime = date.toLocaleTimeString("et-EE");

  return `${formattedDate} ${formattedTime}`;
}

// Convert date to ISO format (yyyy-mm-dd) for backend
export function isoDate(dateString) {
  const [datePart, timePart] = dateString.split(" ");
  const date = datePart.split(".").reverse().join("-");
  const time = timePart.split(":").join(":");

  const parsedDate = new Date(`${date} ${time}`);

  const isoDate = new Date(parsedDate).toISOString();
  console.log("converted datestamp", isoDate);
  return isoDate;
}
