<script lang="ts">
	import type { StateCardType } from '../../types/dashboard';
	export let stateCard: StateCardType;
</script>

<div class="rounded-lg bg-neutral-800 p-4 shadow-lg">
	<div class="flex items-center justify-between">
		<h3 class="text-sm text-neutral-400">{stateCard.label}</h3>
		<svelte:component this={stateCard.icon} />
	</div>
	<p class="mt-2 text-2xl font-bold">{stateCard.value}</p>
</div>
