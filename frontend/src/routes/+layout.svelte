<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { theme } from '$lib/stores/theme';

	let { children } = $props();

	onMount(() => {
		theme.init();
	});
</script>

<svelte:head>
	<title>SyncSnipe — File Sync Engine</title>
	<meta
		name="description"
		content="SyncSnipe: monitor and control your file synchronization rules in real time."
	/>
</svelte:head>

{@render children()}
