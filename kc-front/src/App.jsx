import { useState } from "react";
import "./App.css";

function App() {
  const [formData, setFormData] = useState({
    username: "",
    password: "",
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  async function handleLogin(formData) {
    try {
      const response = await fetch("http://localhost:8081/login", {
        method: "POST",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
          Origin: "http://localhost:5173",
        },
        body: JSON.stringify(formData),
      });

      const jsonData = await response.json();
      if (jsonData) {
        localStorage.setItem("accessToken", jsonData.accessToken);
        localStorage.setItem("refreshToken", jsonData.refreshToken);
      }
    } catch (err) {
      console.log("Error occured during login: ", err);
      return err;
    }
  }

  async function refreshToken(cb) {
    try {
      const refreshToken = localStorage.getItem("refreshToken");

      const response = await fetch("http://localhost:8081/refreshToken", {
        method: "POST",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
          Origin: "http://localhost:5173",
        },
        body: JSON.stringify({ refreshToken }),
      });

      if (response.status === 200) {
        const jsonData = await response.json();
        localStorage.setItem("accessToken", jsonData.accessToken);
        localStorage.setItem("refreshToken", jsonData.refreshToken);
        return cb();
      }
    } catch (err) {
      console.log("Error occured during refresh token request: ", err);
      return err;
    }
  }

  //authorized endpoint
  async function getDocs() {
    try {
      const accessToken = localStorage.getItem("accessToken");

      const response = await fetch("http://localhost:8081/docs", {
        method: "GET",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
          Origin: "http://localhost:5173",
          Authorization: `Bearer ${accessToken}`,
        },
      });

      if (response.status === 401) {
        await refreshToken(getDocs);
      }

      if (response) {
        const jsonData = await response.json();
        const target = document.querySelector("#docs");
        target.textContent = "";

        jsonData.map((el) => {
          const p = document.createElement("p");
          p.textContent = `el.id: ${el.id} el.num:${el.num} el.date:${el.date}`;
          target.appendChild(p);
        });
      }
    } catch (err) {
      console.log("Error occured during getDocs: ", JSON.stringify(err));
      return err;
    }
  }

  async function logout() {
    try {
      console.log("called logout!");
      const accessToken = localStorage.getItem("accessToken");
      const refreshToken = localStorage.getItem("refreshToken");

      const response = await fetch("http://localhost:8081/logout", {
        method: "POST",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
          Origin: "http://localhost:5173",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({ refreshToken }),
      });

      if (response.status === 401) {
        await refreshToken(logout);
      }

      if (response.status === 200) {
        localStorage.removeItem("accessToken");
        localStorage.removeItem("refreshToken");
      }
    } catch (err) {
      console.log("Error occured during logout: ", JSON.stringify(err));
      return err;
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    await handleLogin(formData);
  };

  return (
    <div className="App">
      <h1>Login Form</h1>
      <form onSubmit={handleSubmit}>
        <label>
          Username:
          <input
            type="text"
            name="username"
            value={formData.username}
            onChange={handleChange}
          />
        </label>
        <br />
        <label>
          Password:
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
          />
        </label>
        <br />
        <button type="submit">Login</button>
      </form>

      <button onClick={getDocs}>Show me the docs! (authorized)</button>
      <div id="docs"></div>
      <button onClick={logout}>Logout</button>
    </div>
  );
}

export default App;
