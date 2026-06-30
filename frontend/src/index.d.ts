export type SupportedOperators =
  | "Equal to"
  | "Not equal to"
  | ">"
  | "<"
  | ">="
  | "<="
  | "Like"
  | "Not Like"
  | "In"
  | "Not In"
  | "starts with"
  | "ends with";
export type Condition = {
  operator?: "AND" | "OR";
  column: {
    table: string;
    name: string;
  };
  comparison: SupportedOperators;
  value: {
    payload: null | string;
    value?: string;
    column: {
      table: string;
      name: string;
    };
  };
};
