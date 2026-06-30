export namespace io {
	
	export class Database {
	    from?: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Database(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = source["from"];
	        this.name = source["name"];
	    }
	}
	export class Reference {
	    constraintName: string;
	    sourceTable: string;
	    sourceField: string;
	    localField: string;
	    onDelete: string;
	    onUpdate: string;
	
	    static createFrom(source: any = {}) {
	        return new Reference(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.constraintName = source["constraintName"];
	        this.sourceTable = source["sourceTable"];
	        this.sourceField = source["sourceField"];
	        this.localField = source["localField"];
	        this.onDelete = source["onDelete"];
	        this.onUpdate = source["onUpdate"];
	    }
	}
	export class ForeignKeys {
	    name: string;
	    localField: string;
	    referencedTable: string;
	    referencedField: string;
	    onDelete: string;
	    onUpdate: string;
	
	    static createFrom(source: any = {}) {
	        return new ForeignKeys(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.localField = source["localField"];
	        this.referencedTable = source["referencedTable"];
	        this.referencedField = source["referencedField"];
	        this.onDelete = source["onDelete"];
	        this.onUpdate = source["onUpdate"];
	    }
	}
	export class Index {
	    name: string;
	    column: string;
	    unique?: boolean;
	    indexType: string;
	
	    static createFrom(source: any = {}) {
	        return new Index(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.column = source["column"];
	        this.unique = source["unique"];
	        this.indexType = source["indexType"];
	    }
	}
	export class Field {
	    name: string;
	    type: string;
	    length: number;
	    isPrimaryKey: boolean;
	    nullish: boolean;
	    autoGen: boolean;
	    default: any;
	    attributes: string;
	
	    static createFrom(source: any = {}) {
	        return new Field(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.length = source["length"];
	        this.isPrimaryKey = source["isPrimaryKey"];
	        this.nullish = source["nullish"];
	        this.autoGen = source["autoGen"];
	        this.default = source["default"];
	        this.attributes = source["attributes"];
	    }
	}
	export class Entity {
	    name: string;
	    description: string;
	    primaryKey: string[];
	    onConflict?: string;
	    fields: Field[];
	    indexes?: Index[];
	    ForeignKeys: ForeignKeys[];
	    referencedIn?: Reference[];
	
	    static createFrom(source: any = {}) {
	        return new Entity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.primaryKey = source["primaryKey"];
	        this.onConflict = source["onConflict"];
	        this.fields = this.convertValues(source["fields"], Field);
	        this.indexes = this.convertValues(source["indexes"], Index);
	        this.ForeignKeys = this.convertValues(source["ForeignKeys"], ForeignKeys);
	        this.referencedIn = this.convertValues(source["referencedIn"], Reference);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class ReferencedColumns {
	    constraintName: string;
	    localColumn: string;
	    foreignTable: string;
	    foreignColumn: string;
	    onDelete: string;
	    onUpdate: string;
	
	    static createFrom(source: any = {}) {
	        return new ReferencedColumns(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.constraintName = source["constraintName"];
	        this.localColumn = source["localColumn"];
	        this.foreignTable = source["foreignTable"];
	        this.foreignColumn = source["foreignColumn"];
	        this.onDelete = source["onDelete"];
	        this.onUpdate = source["onUpdate"];
	    }
	}
	export class Schema {
	    entities: Entity[];
	    foreignKeys: ForeignKeys[];
	    referencedColumns: ReferencedColumns[];
	
	    static createFrom(source: any = {}) {
	        return new Schema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], Entity);
	        this.foreignKeys = this.convertValues(source["foreignKeys"], ForeignKeys);
	        this.referencedColumns = this.convertValues(source["referencedColumns"], ReferencedColumns);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Meta {
	    version: string;
	    database: Database;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Meta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.database = this.convertValues(source["database"], Database);
	        this.description = source["description"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class IntermediateSchema {
	    meta: Meta;
	    schema: Schema;
	
	    static createFrom(source: any = {}) {
	        return new IntermediateSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.meta = this.convertValues(source["meta"], Meta);
	        this.schema = this.convertValues(source["schema"], Schema);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ProjectResponse {
	    created: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.created = source["created"];
	        this.message = source["message"];
	    }
	}
	export class ProjectResult {
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	
	

}

