<script lang="ts">
	import type { RecentActivityItemType } from '$lib/types/dashboard';
	export let activity: RecentActivityItemType;
</script>

<!-- Hidden div forces Tailwind to include these color classes in the bundle -->
<div class="hidden bg-green-500 bg-red-500 bg-amber-500 bg-cyan-500 bg-violet-500"></div>

<div class="flex items-start">
	<div class="mt-1.5 mr-3 h-2.5 w-2.5 flex-shrink-0 rounded-full bg-{activity.color}-500"></div>
	<div class="min-w-0 flex-1">
		<p class="text-sm font-medium text-primary leading-snug">{activity.label}</p>
		<p class="text-xs text-secondary mt-0.5 leading-relaxed">{activity.details}</p>
	</div>
</div>