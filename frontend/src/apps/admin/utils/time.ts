export const extractTimeString = (timeVal: any): string => {
  if (!timeVal) return "00:00";
  if (
    typeof timeVal === "string" &&
    timeVal.length <= 5 &&
    timeVal.includes(":")
  ) {
    return timeVal;
  }
  const d = new Date(timeVal);
  if (!isNaN(d.getTime())) {
    const hours = String(d.getHours()).padStart(2, "0");
    const minutes = String(d.getMinutes()).padStart(2, "0");
    return `${hours}:${minutes}`;
  }
  return "00:00";
};

export const transformScheduleToAPI = (data: any) => {
  const cleanTime = extractTimeString(data.start_time_only);
  const startString = `${data.party_day}T${cleanTime}:00`;
  const startDate = new Date(startString);
  const endDate = new Date(startDate.getTime() + data.duration_mins * 60000);

  return {
    ...data,
    start_time: startDate.toISOString(),
    end_time: endDate.toISOString(),
    party_day: undefined,
    start_time_only: undefined,
    duration_mins: undefined,
  };
};
