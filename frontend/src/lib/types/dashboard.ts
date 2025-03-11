import type { Component } from 'svelte';

export type NavItem = {
	href: string;
	text: string;
	icon: Component;
};

export type StateCardType = {
	label: string;
	value: number;
	icon: Component;
};

export type SyncProfileCardType = {
	sourceDir: string;
	targetDir: string;
	status: 'active' | 'ideal' | 'scheduled' | 'paused';
	progress: number;
	lastSync: string;
};

export type RecentActivityItemType = {
	label: string;
	details: string;
	color: string;
}