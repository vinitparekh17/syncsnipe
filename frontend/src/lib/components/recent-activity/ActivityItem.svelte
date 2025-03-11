<script lang="ts">
	import type { RecentActivityItemType } from "$lib/types/dashboard";

    export let activity: RecentActivityItemType;
</script>
<div class="flex items-start">
  <!-- tailwindcss-disable-next-line -->
    <div class="hidden bg-green-500 bg-red-500 bg-amber-500 bg-cyan-500 bg-violet-500"></div>
    <div class="mt-1.5 mr-3 h-3 w-3 flex-shrink-0 rounded-full bg-{activity.color}-500"></div>
    <div>
        <p class="text-sm font-medium">{activity.label}</p>
        <p class="text-xs text-neutral-400">{activity.details}</p>
    </div>
</div>