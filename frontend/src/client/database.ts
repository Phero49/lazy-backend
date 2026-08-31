//not sure why this is named client

import { Notify } from "quasar";
import { useAppContextStore } from "~/store/appContext";
import { DatabaseAttribute } from "..";

const appStore = useAppContextStore();
const buildURl = (path: string, port?: string) => {
  return `http://localhost:${port}${path}`;
};

export function importDatabaseSchema() {}

export function exportDatabaseSchema() {}
export async function importDatabaseColumnAttributes(): Promise<DatabaseAttribute[]> {
  const attributesEndPoints = 
    appStore.activePlugins.database?.plugin["end-points"].attributes;
  
  if (attributesEndPoints?.method.toUpperCase() === "GET") {
    try {
      const response = await fetch(
        buildURl(
          attributesEndPoints.route,
          appStore.activePlugins.database?.port,
        ),
      );
      
      if (!response.ok) {
        const errorMessage = `Failed to get attributes: ${response.status} ${response.statusText}`;
        Notify.create({
          message: errorMessage,
          type: "error"
        });
        throw new Error(errorMessage);
      }
      
      const data: DatabaseAttribute[] = await response.json();
      
      // Optional: Validate the response data
      if (!Array.isArray(data)) {
        throw new Error("Invalid response format: expected array");
      }
      
      return data;
      
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : "Unknown error occurred";
      Notify.create({
        message: `Failed to import database attributes: ${errorMessage}`,
        type: "error"
      });
      throw error;
    }
  } else {
    const errorMessage = `Unexpected HTTP method: ${attributesEndPoints?.method || 'undefined'}. Expected GET`;
    Notify.create({
      message: errorMessage,
      type: "error"
    });
    throw new Error(errorMessage);
  }
}