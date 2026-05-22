const API = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

async function request(path, options = {}) {
  const res = await fetch(`${API}${path}`, {
    headers: { "Content-Type": "application/json", ...(options.headers || {}) },
    ...options,
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export const api = {
  register: (body) => request("/users/register", { method: "POST", body: JSON.stringify(body) }),
  login: (body) => request("/users/login", { method: "POST", body: JSON.stringify(body) }),
  cars: (params = "") => request(`/cars${params}`),
  available: () => request("/cars/available"),
  car: (id) => request(`/cars/${id}`),
  favorite: (carId, UserId) => request(`/cars/${carId}/favorite`, { method: "POST", body: JSON.stringify({ UserId }) }),
  review: (carId, body) => request(`/cars/${carId}/reviews`, { method: "POST", body: JSON.stringify(body) }),
  price: (body) => request("/rentals/price", { method: "POST", body: JSON.stringify(body) }),
  createRental: (body) => request("/rentals", { method: "POST", body: JSON.stringify(body) }),
  myRentals: (id) => request(`/users/${id}/rentals`),
  cancelRental: (id) => request(`/rentals/${id}/cancel`, { method: "POST" }),
  pay: (body) => request("/payments", { method: "POST", body: JSON.stringify(body) }),
};
