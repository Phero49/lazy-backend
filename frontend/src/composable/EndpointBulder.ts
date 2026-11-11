import {
  GetDependencyExports,
  GetDependencyManifest,
  GetProjectUsedDependencies,
} from "#/wailsjs/go/main/App";
import { Dialog, Notify } from "quasar";
import { DependencyMethods } from "..";
import {
  type EndPointList,
  OperationBlockData,
  type Routes,
  useAppContextStore,
} from "../store/appContext";

const appContext = useAppContextStore();

export async function getDependencyAllExports() {
  const exports: DependencyMethods[] = [];
  try {
    const dependencies = await GetProjectUsedDependencies(
      appContext.currentProject!,
    );
    for (const d of dependencies) {
      const m = await GetDependencyManifest(d.location);
      m.modules.forEach((v) => {
        exports.push({
          data: v.exports,
          source: m.location || "",
        });
      });
    }
  } catch (error) {}
  return exports;
}

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
  function getEndPointOperations(index: number) {
    if (route == undefined) {
      return null;
    }
    return route.data?.operations || [];
  }

  const getIndPointOperation = (index: number, blockIndex: number) =>
    route?.data?.operations[index].data[blockIndex];

  const addTable = (table: string, index: number, blockIndex: number) => {
    const opt = getIndPointOperation(index, blockIndex);
    if (opt?.type != "database") {
      return;
    }
    if (opt?.databaseBlock == null) {
      return;
    }
    if (opt.databaseBlock.query.tables == undefined) {
      opt.databaseBlock.query.tables = [];
    }
    opt.databaseBlock.query.tables.push({
      table: table,
    });

    const length = opt.databaseBlock.query.tables.length;
    if (length > 1) {
      const currentEntry = opt.databaseBlock.query.tables[length - 1];
      const previousEntry = opt.databaseBlock.query.tables[length - 2];

      if (currentEntry.join == undefined) {
        currentEntry.join = {
          conditions: [],
          joinType: "INNER",
        };
      }
    }
    appContext.saveApi();
  };
  const getPayload = () => {
    return route?.data?.payload;
  };

  const getTableJoin = ({
    optIndex,
    tableIndex,
    blockIndex,
  }: {
    optIndex: number;
    blockIndex: number;
    tableIndex: number;
  }) => {
    const block = route?.data?.operations[optIndex].data[blockIndex];
    if (block?.type == "database") {
      return block.databaseBlock?.query.tables[tableIndex].join;
    }
  };



function deleteOperationByIndex(blockIndex: number, operationIndex: number) {
  const ops = route?.data?.operations[blockIndex]?.data;
  if (!ops) return;

  // Remove immediately and keep a reference for undo.
  const [removed] = ops.splice(operationIndex, 1);

  Notify.create({
    position:'bottom-right',
    message: "Operation deleted",
    color: "grey-10",
    timeout: 5000,
    actions: [
      {
        label: "Undo",
        color: "yellow",
        handler: () => {
          // Re-insert at the original position to preserve ordering.
          ops.splice(operationIndex, 0, removed);
        },
      },
    ],
  });
}
function deleteBlockByIndex(blockIndex: number) {
  const ops = route?.data?.operations;
  if (!ops || blockIndex < 0 || blockIndex >= ops.length) return;

  const [removed] = ops.splice(blockIndex, 1);

  Notify.create({
    message: "Block deleted",
    color: "grey-9",
    timeout: 5000,
    actions: [{
      label: "Undo",
      color: "yellow",
      handler: () => ops.splice(blockIndex, 0, removed),
    }],
  });
}
  function getEndPointOperationB4Position(operationIndex: number) {
    const index = operationIndex;
    if (route == undefined || index < 0) {
      return [];
    }
    const transformData = (op: OperationBlockData) => {
      if (op.type == "database") {
        return {
          operationIndex: index,
          type: "database",
          name: op.databaseBlock.name,
        };
      } else if (op.type == "dependency") {
        return {
          type: "dependency",
          operationIndex: index,
          name: op.dependenceBlock.name,
        };
      }
    };
    const currentEndPoint = route.data?.operations[index];
    const conditionDisabled = currentEndPoint?.condition?.disableCondition;
    if (!conditionDisabled) {
      const mappedData = route.data?.operations
        .slice(0, index)
        .filter((v) => v.condition?.disableCondition !== false)
        .map((op, i) => {
          return op.data.map(transformData).filter((v) => v != undefined);
        });
      return mappedData?.flat() || [];
    }
    return (
      currentEndPoint?.data.map(transformData).filter((v) => v != undefined) ||
      []
    );
  }
  
  return {
    route,
    getDependencyAllExports,
    getEndPointOperations,
    updateLimitPayload,
    getEndPointOperationB4Position,
    addTable,
    getIndPointOperation,
    deleteOperationByIndex,
    deleteBlockByIndex,
    getPayload,
    getTableJoin,
  };
}

export type EndPointInstance = ReturnType<typeof useEndPointContext>;
