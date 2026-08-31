import { EndPointData } from './../wailsjs/go/models';
import { DatabaseOperations, Routes } from "./store/appContext";

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
 type Condition = {
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


 type UseOperationInfo = () => {
    index: number;
    routeData: Routes;
    dbOperation: DatabaseOperations;
}



interface AttributeParam {
 
    name: string;
    value: 'string' | 'int' | 'float';
    description?: string;
    optional?:boolean,
    variadic?:boolean
  
}

interface DatabaseAttribute {
  name: string;
  description: string;
  type: 'property' | 'function';
  params?: AttributeParam[];  // Renamed from 'param' to 'params' (plural) for clarity
}



type OperationInfo ={
 index:number,
 routeData:EndPointData,
 dbOperation:DatabaseOperations
}