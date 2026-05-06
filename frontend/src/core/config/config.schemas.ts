import { z } from "zod";

const DEFAULT_DURATIONS = [0, 5, 10, 15, 30, 45, 60, 90, 120];

const getDefaultPartyWeekend = () => {
  const today = new Date();
  const dayOfWeek = today.getDay(); // 0 = Sonntag, 1 = Montag ... 5 = Freitag

  const isoDay = dayOfWeek === 0 ? 7 : dayOfWeek;

  const daysToFriday = 5 - isoDay;

  const friday = new Date(today);
  friday.setDate(today.getDate() + daysToFriday);

  const saturday = new Date(friday);
  saturday.setDate(friday.getDate() + 1);

  const sunday = new Date(friday);
  sunday.setDate(friday.getDate() + 2);

  const format = (d: Date) => {
    const yyyy = d.getFullYear();
    const mm = String(d.getMonth() + 1).padStart(2, "0");
    const dd = String(d.getDate()).padStart(2, "0");
    return `${yyyy}-${mm}-${dd}`;
  };

  return [
    { id: format(friday), name: "Friday" },
    { id: format(saturday), name: "Saturday" },
    { id: format(sunday), name: "Sunday" },
  ];
};

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
  date_locale: z.string().default("de-DE"),
  date_options: z.record(z.string(), z.any()).default({
    weekday: "long",
    hour: "2-digit",
    minute: "2-digit",
  }),
  /* @TODO remove the defaults once we have a real API providing these values */
  timezone: z.string().default("Europe/Berlin"),
  party_days: z.array(DaySchema).default(getDefaultPartyWeekend),
  event_durations: DurationsSchema.default(DEFAULT_DURATIONS),
});

export type AppConfig = z.infer<typeof AppConfigSchema>;
export type Day = z.infer<typeof DaySchema>;
export type Durations = z.infer<typeof DurationsSchema>;
