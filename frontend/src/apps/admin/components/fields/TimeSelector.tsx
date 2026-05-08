import { useEffect, useRef } from "react";
import {
  SelectInput,
  TimeInput,
  NumberInput,
  useRecordContext,
  minValue,
  maxValue,
} from "react-admin";
import { useWatch, useFormContext } from "react-hook-form";
import { Box, Typography, Chip, Stack } from "@mui/material";
import EventIcon from "@mui/icons-material/Event";
import { useAppConfig } from "@core/config/useConfig";
import { Day, Durations } from "@core/config/config.schemas";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc);
dayjs.extend(timezone);
interface TimeSelectorProps {
  partyDays: Day[];
  presetDurations?: Durations;
}

export const TimeSelector = ({
  partyDays,
  presetDurations = [15, 30, 60, 90], // Fallback if nothing comes from the backend
}: TimeSelectorProps) => {
  const {
    date_locale: dateLocale,
    date_options: dateOptions,
    timezone: tz,
  } = useAppConfig();

  const record = useRecordContext(); // NEW: Fetches backend data (when in edit view)

  // 1. Hook for the quick buttons
  const { setValue } = useFormContext();

  // 2. Hooks to monitor the three fields
  const partyDay = useWatch({ name: "party_day" });
  const startTime = useWatch({ name: "start_time_only" });
  const duration = useWatch({ name: "duration_mins" });

  // 3. Ref to ensure initialization runs only once
  const initializedRef = useRef(false);

  // 4. Compute preview directly from watched values
  const preview = (() => {
    if (partyDay && startTime && duration !== undefined && duration !== null) {
      const cleanTime =
        String(startTime).length > 5
          ? dayjs(startTime).format("HH:mm")
          : startTime;
      const startObj = dayjs.tz(`${partyDay} ${cleanTime}`, tz);

      if (startObj.isValid()) {
        const endDate = startObj.add(duration, "minute");

        const options: Intl.DateTimeFormatOptions = {
          ...dateOptions,
          timeZone: tz,
        };

        return `Event ends on: ${endDate.toDate().toLocaleString(dateLocale, options)}`;
      }
    }
    return "Please fill out the three fields.";
  })();

  useEffect(() => {
    if (record?.start_time && !initializedRef.current) {
      const startObj = dayjs(record.start_time).tz(tz);

      if (startObj.isValid()) {
        const loadedDay = startObj.format("YYYY-MM-DD");
        const loadedTime = startObj.format("HH:mm");

        let diffMins = 60;
        if (record.end_time) {
          const endObj = dayjs(record.end_time).tz(tz);
          diffMins = endObj.diff(startObj, "minute");
        }

        const timeoutId = setTimeout(() => {
          setValue("party_day", loadedDay, {
            shouldDirty: true,
            shouldValidate: true,
            shouldTouch: true,
          });
          setValue("start_time_only", loadedTime, {
            shouldDirty: true,
            shouldValidate: true,
            shouldTouch: true,
          });
          setValue("duration_mins", diffMins, {
            shouldDirty: true,
            shouldValidate: true,
            shouldTouch: true,
          });
        }, 0);

        initializedRef.current = true;
        return () => {
          clearTimeout(timeoutId);
          initializedRef.current = false;
        };
      }
    }
  }, [record, setValue, tz]);

  return (
    <Box sx={{ mb: 2, p: 2, border: "1px dashed #ccc", borderRadius: 2 }}>
      <Typography variant="subtitle2" color="textSecondary" gutterBottom>
        Schedule Timing
      </Typography>

      <Stack direction="row" spacing={2} alignItems="flex-start">
        <SelectInput
          source="party_day"
          choices={partyDays}
          label="Day"
          required
        />

        <TimeInput source="start_time_only" label="Start Time" required />

        <Box>
          <NumberInput
            source="duration_mins"
            label="Duration (m)"
            required
            min={0}
            max={600}
            sx={{ mt: 0 }}
            validate={[
              minValue(0, "Duration cannot be negative"),
              maxValue(600, "Duration cannot exceed 10 hours"),
            ]}
            onKeyDown={(e) => {
              if (
                e.key === "-" ||
                e.key === "e" ||
                e.key === "+" ||
                e.key === "." ||
                e.key === ","
              ) {
                e.preventDefault();
              }
            }}
          />

          <Stack direction="row" spacing={1} sx={{ mt: -2, mb: 2 }}>
            {presetDurations.map((mins) => {
              const isActive = duration === mins;

              return (
                <Chip
                  key={mins}
                  label={`${mins}m`}
                  size="small"
                  color={isActive ? "primary" : "default"}
                  variant={isActive ? "filled" : "outlined"}
                  onClick={() =>
                    setValue("duration_mins", mins, {
                      shouldDirty: true,
                      shouldValidate: true,
                      shouldTouch: true,
                    })
                  }
                  sx={{
                    cursor: "pointer",
                    fontWeight: isActive ? "bold" : "normal",
                    "&:hover": {
                      backgroundColor: isActive
                        ? "primary.main"
                        : "action.hover",
                    },
                  }}
                />
              );
            })}
          </Stack>
        </Box>
      </Stack>

      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          color: "primary.main",
          mt: 1,
        }}
      >
        <EventIcon fontSize="small" sx={{ mr: 1 }} />
        <Typography variant="body2" fontWeight="bold">
          {preview}
        </Typography>
      </Box>
    </Box>
  );
};
