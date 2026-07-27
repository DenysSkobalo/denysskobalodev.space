CREATE TABLE `content_feed` (
	`id` text PRIMARY KEY NOT NULL,
	`source` text DEFAULT 'telegram' NOT NULL,
	`external_id` text NOT NULL,
	`title` text,
	`payload_json` text NOT NULL,
	`published_at` integer NOT NULL,
	`created_at` integer NOT NULL
);
--> statement-breakpoint
CREATE TABLE `products` (
	`id` text PRIMARY KEY NOT NULL,
	`name` text NOT NULL,
	`description` text,
	`created_at` integer NOT NULL
);
--> statement-breakpoint
CREATE TABLE `sessions` (
	`id` text PRIMARY KEY NOT NULL,
	`user_id` text NOT NULL,
	`token_hash` text NOT NULL,
	`ip_address` text,
	`user_agent` text,
	`expires_at` integer NOT NULL,
	`created_at` integer NOT NULL,
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `user_entitlements` (
	`id` text PRIMARY KEY NOT NULL,
	`email` text NOT NULL,
	`user_id` text,
	`product_id` text NOT NULL,
	`grant_type` text DEFAULT 'paid' NOT NULL,
	`monthly_quota` integer DEFAULT 100 NOT NULL,
	`status` text DEFAULT 'active' NOT NULL,
	`granted_by` text,
	`starts_at` integer NOT NULL,
	`expires_at` integer,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE set null,
	FOREIGN KEY (`product_id`) REFERENCES `products`(`id`) ON UPDATE no action ON DELETE no action
);
--> statement-breakpoint
CREATE TABLE `users` (
	`id` text PRIMARY KEY NOT NULL,
	`email` text NOT NULL,
	`display_name` text NOT NULL,
	`role` text DEFAULT 'user' NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL
);
--> statement-breakpoint
CREATE UNIQUE INDEX `content_feed_external_id_unique` ON `content_feed` (`external_id`);--> statement-breakpoint
CREATE INDEX `idx_content_feed_published` ON `content_feed` (`published_at`);--> statement-breakpoint
CREATE UNIQUE INDEX `sessions_token_hash_unique` ON `sessions` (`token_hash`);--> statement-breakpoint
CREATE INDEX `idx_sessions_token_hash` ON `sessions` (`token_hash`);--> statement-breakpoint
CREATE INDEX `idx_sessions_user_id` ON `sessions` (`user_id`);--> statement-breakpoint
CREATE INDEX `idx_entitlements_email` ON `user_entitlements` (`email`);--> statement-breakpoint
CREATE INDEX `idx_entitlements_user_id` ON `user_entitlements` (`user_id`);--> statement-breakpoint
CREATE INDEX `idx_entitlements_product_status` ON `user_entitlements` (`product_id`,`status`);--> statement-breakpoint
CREATE UNIQUE INDEX `users_email_unique` ON `users` (`email`);
