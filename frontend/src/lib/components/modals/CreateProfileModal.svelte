<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { createProfile } from '$lib/api';

	export let isOpen = false;

	const dispatch = createEventDispatcher();
	let profileName = '';
	let loading = false;
	let error = '';

	async function handleSubmit() {
		error = '';
		if (!profileName.trim()) { error = 'Profile name cannot be empty'; return; }

		loading = true;
		try {
			await createProfile(profileName.trim());
			profileName = '';
			dispatch('created');
			isOpen = false;
		} catch (e: any) {
			error = e.message || 'Failed to create profile';
		} finally {
			loading = false;
		}
	}

	function close() {
		if (loading) return;
		error = '';
		profileName = '';
		isOpen = false;
	}
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4"
		style="background: var(--bg-scrim); backdrop-filter: var(--blur-overlay); -webkit-backdrop-filter: var(--blur-overlay);"
		on:click|self={close}
		on:keydown={(e) => e.key === 'Escape' && close()}
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-create-title"
		tabindex="-1"
	>
		<div class="w-full max-w-md overflow-hidden rounded-2xl border border-themed p-6 surface-raised">
			<!-- Header -->
			<div class="flex items-center justify-between pb-4 border-b border-subtle">
				<h3 id="modal-create-title" class="text-lg font-semibold text-primary" style="letter-spacing: -0.015em;">
					Create New Profile
				</h3>
				<button
					on:click={close}
					class="rounded-full p-1.5 text-secondary hover:bg-inset hover:text-primary transition-colors btn-press"
					aria-label="Close modal"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<form on:submit|preventDefault={handleSubmit} class="mt-4 space-y-4">
				{#if error}
					<div class="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-500">
						{error}
					</div>
				{/if}

				<div>
					<label for="profileName" class="block text-xs font-semibold tracking-widest uppercase text-tertiary mb-1.5">
						Profile Name
					</label>
					<input
						id="profileName"
						type="text"
						bind:value={profileName}
						placeholder="e.g. Work, Photos, Backups"
						disabled={loading}
						class="w-full rounded-xl border border-themed bg-inset px-4 py-2.5 text-sm text-primary
							placeholder-[color:var(--text-tertiary)]
							focus:border-[#8A4FFF] focus:outline-none focus:ring-2 focus:ring-[rgba(138,79,255,0.2)]
							transition-all duration-150 disabled:opacity-50"
					/>
				</div>

				<div class="flex items-center justify-end space-x-3 pt-2">
					<button
						type="button"
						on:click={close}
						disabled={loading}
						class="rounded-xl px-4 py-2 text-sm font-medium text-secondary hover:text-primary transition-colors btn-press"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={loading}
						class="flex items-center rounded-xl bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5]
							px-5 py-2 text-sm font-semibold text-white btn-press hover:opacity-90 disabled:opacity-50"
						style="box-shadow: var(--shadow-brand);"
					>
						{#if loading}
							<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
							</svg>
						{/if}
						Create Profile
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
