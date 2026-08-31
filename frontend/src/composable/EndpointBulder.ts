import { type EndPointList, type Routes, useAppContextStore } from "../store/appContext";

const appContext = useAppContextStore()

export function useEndPointContext(route?: Routes) {
  function updateLimitPayload(v: boolean) {
    if (route == undefined) {
      return null;
    }
    if (v) {
      route.data?.payload.push({
        dataType: "int",
        fieldName: "start_at",
        dataSource: "query",
        required: true,
      });
    } else {
      const index = route.data?.payload.findIndex(
        (v) => v.fieldName == "start_at",
      );
      if (index != undefined && index != -1) {
        route?.data?.payload.splice(index, 1);
      }
    }
  }
  function getEndPointOperations(index:number) {
       if (route == undefined) {
      return null;
    }
   return   route.data?.operations || []
  }

  const getIndPointOperation = (index:number)=>route?.data?.operations[index]

  const addTable = (table: string,index:number) => {
   const opt = getIndPointOperation(index)
   if (opt == null) {
    return
   }
  if (opt.query.tables == undefined) {
    opt.query.tables = [];
  }
  opt.query.tables.push({
    table: table,
  });

  const length = opt.query.tables.length;
  if (length > 1) {
    const currentEntry = opt.query.tables[length - 1];
    const previousEntry = opt.query.tables[length - 2];

    if (currentEntry.join == undefined) {
      currentEntry.join = {
        conditions: [],
        joinType: "INNER",
      };
    }
  }
  appContext.saveApi();
};
const getPayload = ()=>{
    return route?.data?.payload 
}

const getTableJoin = ({optIndex,tableIndex}:{optIndex:number,tableIndex:number})=>{
    console.log( 'operations', route?.data?.operations[optIndex].query.tables)
  return route?.data?.operations[optIndex].query.tables[tableIndex].join
}
  
  return {
    route,
    updateLimitPayload,
    getEndPointOperations,
    addTable,getIndPointOperation,
    getPayload,
    getTableJoin
  };
}

export type EndPointInstance =  ReturnType<typeof useEndPointContext>