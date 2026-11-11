import { EndPointData, io } from './../wailsjs/go/models';
import { ConditionGate, DatabaseOperations, OperationBlock, Routes } from "./store/appContext";

/**
 * Comparison operators supported by the query builder.
 * These describe how a column is compared to another value or column.
 */
type SupportedOperators =
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

/**
 * A single WHERE clause or logical condition used in a database query.
 */
type Condition = {
  /** Logical connector between this condition and the next one. */
  operator?: "AND" | "OR";
  /** Column being compared. */
  column: {
    /** Table name for the source column. */
    table: string;
    /** Column name on the table. */
    name: string;
  };
  /** Comparison operator applied to the column. */
  comparison: SupportedOperators;
  /** Right-hand side of the condition. */
  value: {
    /** Optional payload value used when the source is a request field. */
    payload: null | string;
    /** Literal value when the comparison uses a static string. */
    value?: string;
    /** Column reference when the comparison uses another column. */
    column: {
      /** Source table for the other column. */
      table: string;
      /** Source column name. */
      name: string;
    };
  };
};

/**
 * Utility shape returned by the endpoint operation hook.
 */
type UseOperationInfo = () => {
  /** Index of the current operation inside the route. */
  index: number;
  /** Route data being processed. */
  routeData: Routes;
  /** Database operation associated with this step. */
  dbOperation: DatabaseOperations;
};

/**
 * A parameter definition for a database attribute or helper function.
 */
interface AttributeParam {
  /** Name shown to the user or used in generated code. */
  name: string;
  /** Value type accepted by the param. */
  value: 'string' | 'int' | 'float';
  /** Optional explanation for the parameter. */
  description?: string;
  /** Whether the argument is optional. */
  optional?: boolean;
  /** Whether the parameter can be repeated. */
  variadic?: boolean;
}

/**
 * Metadata for a database attribute definition such as a property or helper function.
 */
interface DatabaseAttribute {
  /** Attribute name exposed to the user. */
  name: string;
  /** Human-readable explanation of what this attribute does. */
  description: string;
  /** Whether this attribute is a property or a method-like function. */
  type: 'property' | 'function';
  /** Optional list of arguments accepted by the attribute. */
  params?: AttributeParam[];
}

/**
 * Detailed information about a database operation nested inside a route.
 */
type OperationInfo = {
  /** Position of the operation within the route. */
  index: number;
  /** Backend endpoint payload for the route being edited. */
  routeData: EndPointData;
  /** Database operation that this entry refers to. */
  dbOperation: DatabaseOperations;
};

type ApiResponse = {
  code:number,
  statusMessage:string,
  userMessage:{
    dataType:string
    message:null|any
  }

}

 type DependencyMethods = {
  source:string,
  data:io.Exports[]
}

type AvailableDataSources = "static" | "module" | "payload" | "operationBlock";

type BlockScope = {
        operation: OperationBlock | undefined;
        condition: ConditionGate | null | undefined;
        enableCondition: boolean;
        operationIndex?:number
        blockIndex:number
        operations: OperationBlock[] | undefined;
        addExpression: () => void;
        hidToggle?:boolean
    };