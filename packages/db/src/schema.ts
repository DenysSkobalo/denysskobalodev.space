import { sqliteTable, text, integer, index } from "drizzle-orm/sqlite-core";

// USERS & IDENTITY
export const users = sqliteTable("users", {
  id: text("id").primaryKey(), // KSUID / ULID string
  email: text("email").notNull().unique(),
  displayName: text("display_name").notNull(),
  role: text("role", { enum: ["admin", "user"] }).notNull().default("user"),
  createdAt: integer("created_at", { mode: "timestamp" }).notNull(),
  updatedAt: integer("updated_at", { mode: "timestamp" }).notNull(),
});

// CROSS-SUBDOMAIN SESSIONS (.denysskobalodev.space)
export const sessions = sqliteTable("sessions", {
  id: text("id").primaryKey(),
  userId: text("user_id").notNull().references(() => users.id, { onDelete: "cascade" }),
  tokenHash: text("token_hash").notNull().unique(),
  ipAddress: text("ip_address"),
  userAgent: text("user_agent"),
  expiresAt: integer("expires_at", { mode: "timestamp" }).notNull(),
  createdAt: integer("created_at", { mode: "timestamp" }).notNull(),
}, (table) => ({
  tokenHashIdx: index("idx_sessions_token_hash").on(table.tokenHash),
  userIdIdx: index("idx_sessions_user_id").on(table.userId),
}));

// PRODUCTS & BETA ENTITLEMENTS
export const products = sqliteTable("products", {
  id: text("id").primaryKey(), 
  name: text("name").notNull(),
  description: text("description"),
  createdAt: integer("created_at", { mode: "timestamp" }).notNull(),
});

export const userEntitlements = sqliteTable("user_entitlements", {
  id: text("id").primaryKey(),
  email: text("email").notNull(), 
  userId: text("user_id").references(() => users.id, { onDelete: "set null" }),
  productId: text("product_id").notNull().references(() => products.id),
  grantType: text("grant_type", { enum: ["paid", "beta_grant", "trial"] }).notNull().default("paid"),
  monthlyQuota: integer("monthly_quota").notNull().default(100), // -1 for unlimited
  status: text("status", { enum: ["active", "suspended", "expired"] }).notNull().default("active"),
  grantedBy: text("granted_by"),
  startsAt: integer("starts_at", { mode: "timestamp" }).notNull(),
  expiresAt: integer("expires_at", { mode: "timestamp" }),
  createdAt: integer("created_at", { mode: "timestamp" }).notNull(),
  updatedAt: integer("updated_at", { mode: "timestamp" }).notNull(),
}, (table) => ({
  emailIdx: index("idx_entitlements_email").on(table.email),
  userIdIdx: index("idx_entitlements_user_id").on(table.userId),
  productStatusIdx: index("idx_entitlements_product_status").on(table.productId, table.status),
}));

// TELEMETRY & CONTENT FEED (TELEGRAM SYNC)
export const contentFeed = sqliteTable("content_feed", {
  id: text("id").primaryKey(),
  source: text("source").notNull().default("telegram"),
  externalId: text("external_id").notNull().unique(),
  title: text("title"),
  payloadJson: text("payload_json").notNull(),
  publishedAt: integer("published_at", { mode: "timestamp" }).notNull(),
  createdAt: integer("created_at", { mode: "timestamp" }).notNull(),
}, (table) => ({
  publishedIdx: index("idx_content_feed_published").on(table.publishedAt),
}));
