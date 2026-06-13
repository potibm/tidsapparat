import { z } from "zod";

const DEFAULT_DURATIONS = [0, 5, 10, 15, 30, 45, 60, 90, 120];

const SentrySchema = z.object({
  dsn: z.string(),
  environment: z.string(),
  version: z.string(),
  replay_session_sample_rate: z.number().min(0).max(1).default(0),
  replay_error_sample_rate: z.number().min(0).max(1).default(1),
});

const DaySchema = z.object({
  id: z.string(),
  name: z.string(),
});

const DurationsSchema = z.array(z.number()).default(DEFAULT_DURATIONS);

export const AppConfigSchema = z.object({
  version: z.string(),
  environment: z.string(),
  environment_message: z.string().optional(),
  sentry: SentrySchema,
  date_locale: z.string().min(2),
  date_options: z.record(z.string(), z.any()).default({
    weekday: "long",
    hour: "2-digit",
    minute: "2-digit",
  }),
  timezone: z.string(),
  party_days: z.array(DaySchema),
  event_durations: DurationsSchema,
  auth: z
    .object({
      type: z.enum(["oidc"]),
      name: z.string(),
      authority: z.string(),
      client_id: z.string(),
    })
    .optional(),
});

export type AppConfig = z.infer<typeof AppConfigSchema>;
export type Day = z.infer<typeof DaySchema>;
export type Durations = z.infer<typeof DurationsSchema>;
