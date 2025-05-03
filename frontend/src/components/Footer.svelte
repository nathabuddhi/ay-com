<script lang="ts">
	// Footer links could be passed as props if needed
	export let showFullFooter: boolean = true;
</script>

<footer>
	<nav>
		{#if showFullFooter}
			<a href="/about">About</a>
			<a href="/download">Download the AY app</a>
			<a href="/help">Help Center</a>
			<a href="/terms">Terms of Service</a>
			<a href="/privacy">Privacy Policy</a>
			<a href="/cookies">Cookie Policy</a>
			<a href="/accessibility">Accessibility</a>
			<a href="/ads">Ads Info</a>
			<a href="/blog">Blog</a>
			<a href="/careers">Careers</a>
			<a href="/resources">Brand Resources</a>
			<a href="/advertising">Advertising</a>
			<a href="/marketing">Marketing</a>
			<a href="/business">AY for Business</a>
			<a href="/developers">Developers</a>
			<a href="/directory">Directory</a>
			<a href="/settings">Settings</a>
		{:else}
			<a href="/terms">Terms</a>
			<a href="/privacy">Privacy</a>
			<a href="/cookies">Cookies</a>
			<a href="/accessibility">Accessibility</a>
		{/if}
		<span>© 2025 AY Corp.</span>
	</nav>
</footer>

<style>
	footer {
		padding: 16px;
		border-top: 1px solid var(--border);
	}

	nav {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 12px;
		font-size: 13px;
		color: var(--secondary);
	}

	a {
		color: var(--secondary);
		text-decoration: none;
	}

	a:hover {
		text-decoration: underline;
	}
</style>
