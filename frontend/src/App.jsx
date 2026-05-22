import React, { useEffect, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import {
  BadgeCheck,
  CalendarDays,
  CarFront,
  CheckCircle2,
  CircleDollarSign,
  Clock3,
  CreditCard,
  Gauge,
  Heart,
  History,
  KeyRound,
  Loader2,
  LogOut,
  MapPin,
  MessageSquareText,
  RefreshCw,
  Search,
  ShieldCheck,
  SlidersHorizontal,
  Sparkles,
  Star,
  UserRound,
  X,
} from "lucide-react";
import { api } from "./lib/api";
import "./styles.css";

const today = new Date();
const tomorrow = new Date(Date.now() + 86400000);
const fallbackImage = "https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&w=1400&q=80";

function App() {
  const [user, setUser] = useState(() => readStoredUser());
  const [cars, setCars] = useState([]);
  const [rentals, setRentals] = useState([]);
  const [selected, setSelected] = useState(null);
  const [favorites, setFavorites] = useState(() => readStoredArray("favorites"));
  const [query, setQuery] = useState({ brand: "", type: "", location: "", availability: "available" });
  const [activeView, setActiveView] = useState("cars");
  const [notice, setNotice] = useState(null);
  const [loadingCars, setLoadingCars] = useState(false);

  const loadCars = async (nextQuery = query) => {
    setLoadingCars(true);
    try {
      const params = new URLSearchParams(Object.fromEntries(Object.entries(nextQuery).filter(([, value]) => value))).toString();
      const data = await api.cars(params ? `?${params}` : "");
      setCars(data.Cars || []);
    } catch (err) {
      showNotice(setNotice, normalizeError(err), "error");
    } finally {
      setLoadingCars(false);
    }
  };

  const loadRentals = async (id = user?.Id) => {
    if (!id) {
      setRentals([]);
      return;
    }
    try {
      const data = await api.myRentals(id);
      setRentals(data.Rentals || []);
    } catch (err) {
      showNotice(setNotice, normalizeError(err), "error");
    }
  };

  useEffect(() => {
    loadCars();
  }, []);

  useEffect(() => {
    loadRentals();
  }, [user?.Id]);

  useEffect(() => {
    localStorage.setItem("favorites", JSON.stringify(favorites));
  }, [favorites]);

  const stats = useMemo(() => {
    const available = cars.filter((car) => car.Status === "available").length;
    const avgRate = cars.length ? Math.round(cars.reduce((sum, car) => sum + Number(car.DailyRate || 0), 0) / cars.length) : 0;
    return { available, avgRate, activeRentals: rentals.filter((r) => ["confirmed", "active", "pending"].includes(r.Status)).length };
  }, [cars, rentals]);

  const visibleCars = useMemo(() => {
    if (activeView !== "favorites") return cars;
    return cars.filter((car) => favorites.includes(car.Id));
  }, [activeView, cars, favorites]);

  const saveSession = (data) => {
    setUser(data.User);
    localStorage.setItem("user", JSON.stringify(data.User));
    localStorage.setItem("token", data.Token || "");
    showNotice(setNotice, `Welcome, ${data.User?.Name || "driver"}.`, "success");
  };

  const signOut = () => {
    setUser(null);
    setRentals([]);
    localStorage.removeItem("user");
    localStorage.removeItem("token");
    showNotice(setNotice, "Signed out.", "success");
  };

  const toggleFavorite = async (car) => {
    if (!user) {
      showNotice(setNotice, "Login first to save favorites.", "error");
      return;
    }
    setFavorites((items) => items.includes(car.Id) ? items.filter((id) => id !== car.Id) : [...items, car.Id]);
    try {
      await api.favorite(car.Id, user.Id);
      showNotice(setNotice, `${car.Brand} ${car.Model} saved to favorites.`, "success");
    } catch {
      showNotice(setNotice, "Saved locally. Backend favorite already exists or could not be updated.", "info");
    }
  };

  return (
    <main className="app-shell">
      <header className="app-header">
        <div className="brand-lockup">
          <span className="brand-mark"><CarFront size={27} /></span>
          <div>
            <strong>DriveFleet</strong>
            <span>Car rental operations</span>
          </div>
        </div>
        <div className="header-actions">
          <button className={activeView === "cars" ? "nav-button active" : "nav-button"} onClick={() => setActiveView("cars")}>
            <CarFront size={18} /> Cars
          </button>
          <button className={activeView === "rentals" ? "nav-button active" : "nav-button"} onClick={() => setActiveView("rentals")}>
            <History size={18} /> Rentals
          </button>
          <button className={activeView === "favorites" ? "nav-button active" : "nav-button"} onClick={() => setActiveView("favorites")}>
            <Heart size={18} /> Favorites
          </button>
          {user ? (
            <button className="user-pill" onClick={signOut} title="Sign out">
              <UserRound size={18} /> {user.Name} <LogOut size={16} />
            </button>
          ) : null}
        </div>
      </header>

      {notice ? <Toast notice={notice} onClose={() => setNotice(null)} /> : null}

      <section className="hero-band">
        <div className="hero-copy">
          <div className="eyebrow"><Sparkles size={16} /> Live fleet booking</div>
          <h1>Reserve the right car, track payments, and keep rentals moving.</h1>
          <p>Connected to Go microservices through the API Gateway with PostgreSQL data, Redis availability caching, NATS events, and email notifications.</p>
        </div>
        <div className="hero-stats">
          <Metric icon={<BadgeCheck />} label="Available cars" value={stats.available} />
          <Metric icon={<CircleDollarSign />} label="Average daily rate" value={`$${stats.avgRate}`} />
          <Metric icon={<Clock3 />} label="Active rentals" value={stats.activeRentals} />
        </div>
      </section>

      {!user ? <AuthCard onLogin={saveSession} setNotice={setNotice} /> : null}

      <section className="main-layout">
        <aside className="control-rail">
          <div className="rail-title"><SlidersHorizontal size={18} /> Search filters</div>
          <Field label="Brand" value={query.brand} onChange={(value) => setQuery({ ...query, brand: value })} placeholder="Toyota, BMW..." />
          <Field label="Type" value={query.type} onChange={(value) => setQuery({ ...query, type: value })} placeholder="sedan, suv" />
          <Field label="Location" value={query.location} onChange={(value) => setQuery({ ...query, location: value })} placeholder="Almaty" />
          <label className="field">
            <span>Availability</span>
            <select value={query.availability} onChange={(e) => setQuery({ ...query, availability: e.target.value })}>
              <option value="available">Available</option>
              <option value="">All cars</option>
              <option value="rented">Rented</option>
            </select>
          </label>
          <button className="primary-button stretch" onClick={() => loadCars()}>
            {loadingCars ? <Loader2 className="spin" size={18} /> : <Search size={18} />} Search fleet
          </button>
          <button className="quiet-button stretch" onClick={() => {
            const clean = { brand: "", type: "", location: "", availability: "available" };
            setQuery(clean);
            loadCars(clean);
          }}>
            <RefreshCw size={18} /> Reset
          </button>
        </aside>

        <section className="content-area">
          {activeView === "rentals" ? (
            <RentalsBoard rentals={rentals} onCancel={async (id) => {
              await api.cancelRental(id);
              await loadRentals();
              await loadCars();
              showNotice(setNotice, "Rental cancelled and car returned to available fleet.", "success");
            }} />
          ) : (
            <>
              <div className="section-heading">
                <div>
                  <h2>{activeView === "favorites" ? "Favorite cars" : "Available fleet"}</h2>
                  <p>{visibleCars.length} cars match the current view</p>
                </div>
                <button className="quiet-button" onClick={() => loadCars()}><RefreshCw size={18} /> Refresh</button>
              </div>
              <div className="car-grid">
                {visibleCars.map((car) => (
                  <CarCard
                    key={car.Id}
                    car={car}
                    favorite={favorites.includes(car.Id)}
                    onFavorite={() => toggleFavorite(car)}
                    onSelect={() => setSelected(car)}
                  />
                ))}
              </div>
              {!visibleCars.length ? <EmptyState title="No cars found" text="Change filters or switch back to all cars." /> : null}
            </>
          )}
        </section>
      </section>

      {selected ? (
        <CarDrawer
          car={selected}
          user={user}
          favorite={favorites.includes(selected.Id)}
          onClose={() => setSelected(null)}
          onFavorite={() => toggleFavorite(selected)}
          onDone={async (message) => {
            setSelected(null);
            await loadRentals();
            await loadCars();
            showNotice(setNotice, message, "success");
          }}
          setNotice={setNotice}
        />
      ) : null}
    </main>
  );
}

function AuthCard({ onLogin, setNotice }) {
  const [mode, setMode] = useState("register");
  const [busy, setBusy] = useState(false);
  const [form, setForm] = useState(() => ({
    Name: "Demo Customer",
    Email: `demo${Date.now()}@example.com`,
    Password: "password123",
    Role: "customer",
  }));

  const submit = async (event) => {
    event.preventDefault();
    setBusy(true);
    try {
      if (mode === "register") {
        await api.register(form);
      }
      const session = await api.login({ Email: form.Email, Password: form.Password });
      onLogin(session);
    } catch (err) {
      showNotice(setNotice, normalizeError(err), "error");
    } finally {
      setBusy(false);
    }
  };

  return (
    <section className="auth-panel">
      <div className="auth-copy">
        <KeyRound size={22} />
        <div>
          <h2>{mode === "register" ? "Create a customer profile" : "Welcome back"}</h2>
          <p>Use the generated demo email for a clean registration, or switch to login for an existing account.</p>
        </div>
      </div>
      <form onSubmit={submit} className="auth-form">
        <div className="segmented">
          <button type="button" className={mode === "register" ? "active" : ""} onClick={() => setMode("register")}>Register</button>
          <button type="button" className={mode === "login" ? "active" : ""} onClick={() => setMode("login")}>Login</button>
        </div>
        {mode === "register" ? <Field label="Name" value={form.Name} onChange={(Name) => setForm({ ...form, Name })} /> : null}
        <Field label="Email" value={form.Email} onChange={(Email) => setForm({ ...form, Email })} />
        <Field label="Password" type="password" value={form.Password} onChange={(Password) => setForm({ ...form, Password })} />
        <button className="primary-button" disabled={busy}>
          {busy ? <Loader2 className="spin" size={18} /> : <UserRound size={18} />} Continue
        </button>
      </form>
    </section>
  );
}

function CarCard({ car, favorite, onFavorite, onSelect }) {
  return (
    <article className="fleet-card">
      <div className="image-wrap">
        <img src={car.ImageUrl || fallbackImage} alt={`${car.Brand} ${car.Model}`} />
        <button className={favorite ? "favorite active" : "favorite"} onClick={onFavorite} title="Favorite">
          <Heart size={18} />
        </button>
        <span className={car.Status === "available" ? "status available" : "status rented"}>{car.Status}</span>
      </div>
      <div className="fleet-body">
        <div>
          <h3>{car.Brand} {car.Model}</h3>
          <p><MapPin size={15} /> {car.Location} <span>{car.Type}</span></p>
        </div>
        <div className="card-metrics">
          <span><CircleDollarSign size={16} /> ${Number(car.DailyRate || 0).toFixed(0)}/day</span>
          <span><Star size={16} /> {Number(car.Rating || 0).toFixed(1)}</span>
        </div>
        <button className="primary-button stretch" onClick={onSelect} disabled={car.Status !== "available"}>
          <CalendarDays size={18} /> {car.Status === "available" ? "Rent this car" : "Unavailable"}
        </button>
      </div>
    </article>
  );
}

function CarDrawer({ car, user, favorite, onClose, onFavorite, onDone, setNotice }) {
  const [form, setForm] = useState({
    UserId: user?.Id || "",
    CarId: car.Id,
    StartDate: toISODate(today),
    EndDate: toISODate(tomorrow),
    Insurance: true,
  });
  const [price, setPrice] = useState(null);
  const [busy, setBusy] = useState(false);
  const [review, setReview] = useState({ Rating: 5, Comment: "" });

  useEffect(() => {
    let alive = true;
    api.price(form).then((data) => alive && setPrice(data)).catch(() => alive && setPrice(null));
    return () => { alive = false; };
  }, [form]);

  const book = async () => {
    if (!user) {
      showNotice(setNotice, "Login or register before booking.", "error");
      return;
    }
    setBusy(true);
    try {
      const rental = await api.createRental({ ...form, UserId: user.Id });
      await api.pay({ RentalId: rental.Id, UserId: user.Id, Amount: rental.TotalPrice, Method: "card" });
      await onDone("Rental confirmed, payment succeeded, and notification event was published.");
    } catch (err) {
      showNotice(setNotice, normalizeError(err), "error");
    } finally {
      setBusy(false);
    }
  };

  const submitReview = async () => {
    if (!user) {
      showNotice(setNotice, "Login first to leave a review.", "error");
      return;
    }
    try {
      await api.review(car.Id, { UserId: user.Id, Rating: Number(review.Rating), Comment: review.Comment });
      showNotice(setNotice, "Review saved. Refresh fleet to update the rating.", "success");
      setReview({ Rating: 5, Comment: "" });
    } catch (err) {
      showNotice(setNotice, normalizeError(err), "error");
    }
  };

  return (
    <div className="drawer-backdrop" role="presentation">
      <aside className="car-drawer">
        <button className="icon-button close-drawer" onClick={onClose} title="Close"><X size={20} /></button>
        <img className="drawer-image" src={car.ImageUrl || fallbackImage} alt={`${car.Brand} ${car.Model}`} />
        <div className="drawer-content">
          <div className="drawer-title">
            <div>
              <span className="eyebrow"><Gauge size={15} /> {car.Type}</span>
              <h2>{car.Brand} {car.Model}</h2>
              <p><MapPin size={16} /> {car.Location}</p>
            </div>
            <button className={favorite ? "icon-button favorite-inline active" : "icon-button favorite-inline"} onClick={onFavorite} title="Favorite">
              <Heart size={20} />
            </button>
          </div>

          <div className="booking-card">
            <h3>Booking details</h3>
            <div className="booking-grid">
              <Field label="Start date" type="date" value={form.StartDate} onChange={(StartDate) => setForm({ ...form, StartDate })} />
              <Field label="End date" type="date" value={form.EndDate} onChange={(EndDate) => setForm({ ...form, EndDate })} />
            </div>
            <label className="insurance-toggle">
              <input type="checkbox" checked={form.Insurance} onChange={(event) => setForm({ ...form, Insurance: event.target.checked })} />
              <span><ShieldCheck size={18} /> Add insurance coverage</span>
            </label>
            <div className="price-strip">
              <span>{price ? `${price.Days} rental days` : "Calculating"}</span>
              <strong>{price ? `$${Number(price.Total || 0).toFixed(2)}` : "..."}</strong>
            </div>
            <button className="primary-button stretch" disabled={busy || car.Status !== "available"} onClick={book}>
              {busy ? <Loader2 className="spin" size={18} /> : <CreditCard size={18} />} Book and pay
            </button>
          </div>

          <div className="review-card">
            <h3>Leave a review</h3>
            <div className="review-row">
              <select value={review.Rating} onChange={(event) => setReview({ ...review, Rating: event.target.value })}>
                {[5, 4, 3, 2, 1].map((rating) => <option key={rating} value={rating}>{rating} stars</option>)}
              </select>
              <button className="quiet-button" onClick={submitReview}><MessageSquareText size={18} /> Save review</button>
            </div>
            <textarea value={review.Comment} onChange={(event) => setReview({ ...review, Comment: event.target.value })} placeholder="Smooth ride, clean cabin, easy pickup..." />
          </div>
        </div>
      </aside>
    </div>
  );
}

function RentalsBoard({ rentals, onCancel }) {
  if (!rentals.length) {
    return <EmptyState title="No rentals yet" text="Book an available car and it will appear here with status and payment history." />;
  }
  return (
    <section className="rentals-board">
      <div className="section-heading">
        <div>
          <h2>My rentals</h2>
          <p>Track current bookings and cancel upcoming rentals.</p>
        </div>
      </div>
      {rentals.map((rental) => (
        <article className="rental-row" key={rental.Id}>
          <div className="rental-icon"><CalendarDays size={20} /></div>
          <div>
            <strong>{rental.StartDate} to {rental.EndDate}</strong>
            <span>Car ID {rental.CarId.slice(0, 8)} · Rental {rental.Id.slice(0, 8)}</span>
          </div>
          <span className={`rental-status ${rental.Status}`}>{rental.Status}</span>
          <strong>${Number(rental.TotalPrice || 0).toFixed(2)}</strong>
          {["confirmed", "active", "pending"].includes(rental.Status) ? (
            <button className="quiet-button danger" onClick={() => onCancel(rental.Id)}>Cancel</button>
          ) : null}
        </article>
      ))}
    </section>
  );
}

function Metric({ icon, label, value }) {
  return (
    <div className="metric">
      <span>{icon}</span>
      <div>
        <strong>{value}</strong>
        <p>{label}</p>
      </div>
    </div>
  );
}

function Field({ label, value, onChange, type = "text", placeholder = "" }) {
  return (
    <label className="field">
      <span>{label}</span>
      <input type={type} value={value} placeholder={placeholder} onChange={(event) => onChange(event.target.value)} />
    </label>
  );
}

function Toast({ notice, onClose }) {
  return (
    <div className={`toast ${notice.kind}`}>
      {notice.kind === "success" ? <CheckCircle2 size={19} /> : <ShieldCheck size={19} />}
      <span>{notice.text}</span>
      <button onClick={onClose} title="Close"><X size={16} /></button>
    </div>
  );
}

function EmptyState({ title, text }) {
  return (
    <div className="empty-state">
      <CarFront size={30} />
      <h3>{title}</h3>
      <p>{text}</p>
    </div>
  );
}

function readStoredUser() {
  try {
    return JSON.parse(localStorage.getItem("user") || "null");
  } catch {
    localStorage.removeItem("user");
    localStorage.removeItem("token");
    return null;
  }
}

function readStoredArray(key) {
  try {
    const value = JSON.parse(localStorage.getItem(key) || "[]");
    return Array.isArray(value) ? value : [];
  } catch {
    return [];
  }
}

function showNotice(setNotice, text, kind = "info") {
  setNotice({ text, kind });
}

function normalizeError(err) {
  const text = String(err?.message || err || "Something went wrong").trim();
  if (text.includes("users_email_key") || text.includes("duplicate key")) return "Email already exists. Login with this email or use another one.";
  if (text.includes("car is not available")) return "This car is no longer available. Refresh the fleet and choose another one.";
  if (text.includes("invalid email or password")) return "Invalid email or password.";
  return text.replace(/^rpc error: code = Unknown desc = /, "");
}

function toISODate(date) {
  return date.toISOString().slice(0, 10);
}

const root = document.getElementById("root");
if (root) {
  createRoot(root).render(<App />);
}
