import { mount } from "svelte";
import "./app.css";
import HomePage from "./pages/HomePage.svelte";

const app = mount(HomePage, {
    target: document.getElementById("app")!,
});

export default app;


