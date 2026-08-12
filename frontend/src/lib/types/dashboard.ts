import type { Component } from 'svelte';
import type { SyncRule } from '$lib/api';

export type NavItem = {
	href: string;
	text: string;
	icon: Component;
};

export type StateCardType = {
	label: string;
	value: string | number;
	icon: Component;
};

export type SyncProfileCardType = {
	id: number;
	name: string;
	createdAt: number;
	updatedAt: number;
	syncRules: SyncRule[];
};

export type RecentActivityItemType = {
	label: string;
	details: string;
	color: string;
};