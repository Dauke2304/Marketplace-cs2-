
let priceSortAsc = true;
let nameSortAsc = true;
let cartItems = [];

function searchItems() {
  const input = document.getElementById("searchInput").value.toLowerCase();
  const cards = document.getElementsByClassName("item-card");

  Array.from(cards).forEach(function (card) {
    const title = card.querySelector(".card-title").textContent.toLowerCase();
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
    return nameSortAsc ? nameA.localeCompare(nameB) : nameB.localeCompare(nameA);
  });

  items.forEach(item => container.appendChild(item));
  nameSortAsc = !nameSortAsc;
}

function openModal(itemName, itemPrice, itemImage) {
  document.getElementById("modalTitle").textContent = itemName;
  document.getElementById("modalPrice").textContent = "Price: $" + itemPrice;
  document.getElementById("modalImage").src = itemImage;
  document.getElementById("cartItemName").value = itemName;
  document.getElementById("cartItemPrice").value = itemPrice;
  const myModal = new bootstrap.Modal(document.getElementById("itemModal"), {});
  myModal.show();
}

function addToCart() {
  const itemName = document.getElementById("cartItemName").value;
  const itemPrice = parseFloat(document.getElementById("cartItemPrice").value);
  cartItems.push({ name: itemName, price: itemPrice });
  updateCart();
}

function updateCart() {
  const totalItems = cartItems.length;
  const cartButton = document.getElementById("cartButton");
  cartButton.innerHTML = `<i class="fas fa-shopping-cart"></i> ${totalItems}`;

  const totalPrice = cartItems.reduce((total, item) => total + item.price, 0);

  if (totalItems > 0) {
    const cartModalBody = document.getElementById("cartModalBody");
    cartModalBody.innerHTML = "";
    cartItems.forEach((item, index) => {
      const itemElement = document.createElement("div");
      itemElement.classList.add("d-flex", "justify-content-between", "align-items-center", "my-2");

      const itemInfo = document.createElement("span");
      itemInfo.textContent = `${item.name} - $${item.price.toFixed(2)}`;

      const removeButton = document.createElement("button");
      removeButton.classList.add("btn", "btn-danger", "btn-sm");
      removeButton.textContent = "Remove";
      removeButton.onclick = () => removeFromCart(index);

      itemElement.appendChild(itemInfo);
      itemElement.appendChild(removeButton);
      cartModalBody.appendChild(itemElement);
    });

    const totalPriceElement = document.createElement("div");
    totalPriceElement.classList.add("fw-bold", "mt-3");
    totalPriceElement.textContent = `Total Price: $${totalPrice.toFixed(2)}`;
    cartModalBody.appendChild(totalPriceElement);
  }
}

function removeFromCart(index) {
  cartItems.splice(index, 1);
  updateCart();
}
