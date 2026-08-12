<script lang="ts">
	import type { StateCardType } from '$lib/types/dashboard';
	export let stateCard: StateCardType;

	// Ensure numeric 0 renders as the string "0" rather than being falsy
	$: displayValue = stateCard.value === 0 || stateCard.value === '0'
		? '0'
		: (stateCard.value ?? '—');
</script>

<div class="group overflow-hidden rounded-2xl border border-themed p-5 surface card-lift">
	<div class="flex items-center justify-between">
		<h3 class="text-xs font-semibold tracking-widest uppercase text-tertiary">
			{stateCard.label}
		</h3>
		<div class="rounded-xl bg-inset p-2 text-[#8A4FFF] group-hover:text-[#20E6B5] transition-colors duration-200">
			<svelte:component this={stateCard.icon} />
		</div>
	</div>
	<p class="mt-3 text-3xl font-bold text-primary" style="letter-spacing: -0.02em; line-height: 1;">
		{displayValue}
	</p>
</div>
