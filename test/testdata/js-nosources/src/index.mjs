import * as css from "./test.css";

console.log(css);


function component() {
  const element = document.createElement('div');

  element.innerHTML = "Hello webpack!";
  element.innerHTML += style;

  return element;
}

document.body.appendChild(component());