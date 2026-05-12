import {
  useRecordContext,
  useResourceContext,
  useUpdate,
  useNotify,
  FieldProps,
} from "react-admin";
import { Switch } from "@mui/material";

export interface BooleanToggleFieldProps extends FieldProps {
  source: string;
}

export const BooleanToggleField = ({
  source,
  ...props
}: BooleanToggleFieldProps) => {
  const record = useRecordContext<Record<string, unknown>>();
  const resource = useResourceContext();
  const notify = useNotify();

  const [update, { isLoading }] = useUpdate();

  if (!record) return null;

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    event.stopPropagation();

    const newValue = event.target.checked;

    update(
      resource,
      {
        id: record.id,
        data: { ...record, [source]: newValue },
        previousData: record,
      },
      {
        mutationMode: "optimistic",
        onSuccess: () => {
          notify(`Status updated`, { type: "success" });
        },
        onError: (error) => {
          const errorMessage =
            error instanceof Error ? error.message : String(error);

          notify(`Error while updating status: ${errorMessage}`, {
            type: "error",
          });
        },
      },
    );
  };

  return (
    <Switch
      checked={!!record[source]}
      onChange={handleChange}
      disabled={isLoading}
      onClick={(e) => e.stopPropagation()}
      {...props}
    />
  );
};

BooleanToggleField.defaultProps = {
  addLabel: true,
};
