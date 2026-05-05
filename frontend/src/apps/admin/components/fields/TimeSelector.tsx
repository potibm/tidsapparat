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
import { extractTimeString } from "@admin/utils/time";

// Interfaces für the configuration
interface PartyDay {
  id: string; // z.B. "2026-05-08"
  name: string; // z.B. "Freitag"
}

interface TimeSelectorProps {
  partyDays: PartyDay[];
  presetDurations?: number[]; // z.B. [15, 30, 45, 60, 90, 120]
}

export const TimeSelector = ({
  partyDays,
  presetDurations = [15, 30, 60, 90], // Fallback, falls vom Backend nichts kommt
}: TimeSelectorProps) => {
  const record = useRecordContext(); // NEU: Holt sich die Backend-Daten (falls im Edit-View)

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
      const cleanTime = extractTimeString(startTime);
      const startStr = `${partyDay}T${cleanTime}:00`;
      const startDate = new Date(startStr);

      if (!isNaN(startDate.getTime())) {
        const endDate = new Date(startDate.getTime() + duration * 60000);

        const options: Intl.DateTimeFormatOptions = {
          weekday: "short",
          hour: "2-digit",
          minute: "2-digit",
        };
        return `Event ends on: ${endDate.toLocaleString("en-GB", options)}`;
      }
    }
    return "Please fill out the three fields.";
  })();

  useEffect(() => {
    if (record && record.start_time && !initializedRef.current) {
      const startDate = new Date(record.start_time);

      if (!isNaN(startDate.getTime())) {
        const year = startDate.getFullYear();
        const month = String(startDate.getMonth() + 1).padStart(2, "0");
        const day = String(startDate.getDate()).padStart(2, "0");

        const loadedDay = `${year}-${month}-${day}`;
        const loadedTime = extractTimeString(startDate);

        let diffMins = 60;
        if (record.end_time) {
          const endDate = new Date(record.end_time);
          diffMins = Math.round(
            (endDate.getTime() - startDate.getTime()) / 60000,
          );
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
        return () => clearTimeout(timeoutId);
      }

      initializedRef.current = true;
    }
  }, [record, setValue]);

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
