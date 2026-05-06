import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc);
dayjs.extend(timezone);

export const transformScheduleToAPI = (
  data: {
    start_time_only: string | Date | number;
    party_day: string;
    duration_mins: number;
    [key: string]: unknown;
  },
  tz: string,
) => {
  const rawTime = String(data.start_time_only);
  const cleanTime =
    rawTime.length > 5 ? dayjs(rawTime).format("HH:mm") : rawTime;

  const startObj = dayjs.tz(`${data.party_day} ${cleanTime}`, tz);

  const endObj = startObj.add(data.duration_mins, "minute");

  return {
    ...data,
    start_time: startObj.toISOString(),
    end_time: endObj.toISOString(),

    party_day: undefined,
    start_time_only: undefined,
    duration_mins: undefined,
  };
};
