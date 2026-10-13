import { io } from "#/wailsjs/go/models";
import router from "~/routes";
import { useAppContextStore } from "~/store/appContext";

const appContext = useAppContextStore();

/**
 * Finds every foreign key in the schema that points at `currentTableName`
 * and returns them as reverse relationships ("referenced in").
 *
 * @param entities         All tables in the schema.
 * @param currentTableName The table whose incoming references we want.
 */
export function computeReferencedIn(
  entities: io.Entity[],
  currentTableName: string,
): NonNullable<io.Entity["referencedIn"]> {
  const referencedIn: NonNullable<io.Entity["referencedIn"]> = [];

  for (const sourceEntity of entities) {
    // `e` owns these foreign keys, so it is the source table.
    for (const fk of sourceEntity.ForeignKeys ?? []) {
      // Only keep foreign keys that point to the current table.
      if (fk.referencedTable !== currentTableName) continue;

      referencedIn.push({
        // Name of the foreign-key constraint.
        constraintName: fk.name,
        // Field on the current (referenced) table.
        localField: fk.referencedField,
        // Field on the table that owns the foreign key.
        sourceField: fk.localField,
        // Table that owns the foreign key.
        sourceTable: sourceEntity.name,
        // Actions when the referenced row is deleted / updated.
        onDelete: fk.onDelete,
        onUpdate: fk.onUpdate,
      });
    }
  }

  return referencedIn;
}

/**
 * Opens or focuses the given tab and highlights the target column.
 */
export function viewColumnColumnSource(tabPath: string,tableName:string, column: string) {
  const safePath = tabPath?.trim();
  const safeColumn = column?.trim();

  if (!safePath || !safeColumn) {
    return;
  }

  const normalizedPath = safePath.startsWith("/") ? safePath : `/${safePath}`;
  const existingTab = appContext.openedTabs.find((tab) => {
    const tabPathNoQuery = tab.path.split("?")[0];
    return (
      tabPathNoQuery === normalizedPath ||
      tabPathNoQuery.includes(normalizedPath) ||
      normalizedPath.includes(tabPathNoQuery)
    );
  });

  const query = {
    ...(existingTab?.query ?? {}),
    highlightColumn: safeColumn,
  };

  if (existingTab) {
    existingTab.query = query;
    appContext.setActiveTab(existingTab.id);
    router.replace({
      path: existingTab.path.split("?")[0],
      query,
    });
    return;
  }
openNewDatabaseTab(normalizedPath, tableName, query);
}


export function openNewDatabaseTab(tabPath: string,label:string, query?: Record<string, string> ) {
  appContext.createNewTab(tabPath, label,'bi-table', query);

}

  