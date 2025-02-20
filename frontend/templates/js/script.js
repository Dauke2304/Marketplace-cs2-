
let priceSortAsc = true;
let nameSortAsc = true;
let cartItems = [];

function searchItems() {
  const input = document.getElementById("searchInput").value.toLowerCase();
  const cards = document.getElementsByClassName("item-card");

  Array.from(cards).forEach(function (card) {
    const title = card.querySelector("h5").textContent.toLowerCase();
    card.style.display = title.includes(input) ? "block" : "none";
  });
}

function sortItemsByPrice() {
  const container = document.getElementById("item-container");
  const items = Array.from(container.getElementsByClassName("item-card"));

  items.sort((a, b) => {
    const priceA = parseFloat(a.getAttribute("data-price"));
    const priceB = parseFloat(b.getAttribute("data-price"));
    return priceSortAsc ? priceA - priceB : priceB - priceA;
  });

  items.forEach(item => container.appendChild(item));
  priceSortAsc = !priceSortAsc;
}

function sortItemsByName() {
  const container = document.getElementById("item-container");
  const items = Array.from(container.getElementsByClassName("item-card"));

  items.sort((a, b) => {
      const nameA = a.querySelector(".card-title").textContent.toLowerCase();
      const nameB = b.querySelector(".card-title").textContent.toLowerCase();
      return nameA.localeCompare(nameB);
  });

  container.innerHTML = "";
  items.forEach(item => container.appendChild(item));
}
