<script lang="ts">
    import { onMount } from "svelte";
    import { initTheme } from "./stores/theme-wrapper";
    import "./styles/app.scss";
    import { Router, Route } from "svelte-routing";
    import NotFoundPage from "./pages/NotFoundPage.svelte";
    import LandingPage from "./pages/LandingPage.svelte";
    import LoginPage from "./pages/LoginPage.svelte";
    import RegisterPage from "./pages/RegisterPage.svelte";
    import VerificationPage from "./pages/VerificationPage.svelte";
    import ForgotPasswordPage from "./pages/ForgotPasswordPage.svelte";
    import ResetPasswordPage from "./pages/ResetPasswordPage.svelte";
    import HomePage from "./pages/HomePage.svelte";

    let route: string = "/";

    function handleRouteChange(): void {
        route = window.location.pathname;
    }

    onMount(() => {
        initTheme();

        window.addEventListener("popstate", handleRouteChange);
        handleRouteChange();

        return () => {
            window.removeEventListener("popstate", handleRouteChange);
        };
    });
</script>

<svelte:head>
    <title>AY.com - Connect, share, engage</title>
</svelte:head>

<Router>
    <Route path="/">
        <LandingPage />
    </Route>
    <Route path="/login">
        <LoginPage />
    </Route>
    <Route path="/register">
        <RegisterPage />
    </Route>
    <Route path="/home">
        <HomePage />
    </Route>
    <Route path="/verification">
        <VerificationPage />
    </Route>
    <Route path="/forgot">
        <ForgotPasswordPage />
    </Route>
    <Route path="/reset-password">
        <ResetPasswordPage />
    </Route>
    <Route path="*">
        <NotFoundPage />
    </Route>
</Router>
