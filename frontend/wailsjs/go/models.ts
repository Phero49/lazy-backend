export namespace io {
	
	export class Index {
	    name: string;
	    fields: string[];
	    unique?: boolean;
	    sparse?: boolean;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Index(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.fields = source["fields"];
	        this.unique = source["unique"];
	        this.sparse = source["sparse"];
	        this.description = source["description"];
	    }
	}
	export class Field {
	    name: string;
	    type: string;
	    subtype?: string;
	    length?: number;
	    nulish: boolean;
	    default?: any;
	    embed?: boolean;
	    description: string;
	    enumValues?: string[];
	    fields?: Field[];
	    items?: Field;
	    valueType?: string;
	
	    static createFrom(source: any = {}) {
	        return new Field(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.subtype = source["subtype"];
	        this.length = source["length"];
	        this.nulish = source["nulish"];
	        this.default = source["default"];
	        this.embed = source["embed"];
	        this.description = source["description"];
	        this.enumValues = source["enumValues"];
	        this.fields = this.convertValues(source["fields"], Field);
	        this.items = this.convertValues(source["items"], Field);
	        this.valueType = source["valueType"];
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
	export class Entity {
	    name: string;
	    description: string;
	    primaryKey: string[];
	    onConflict?: string;
	    fields: Field[];
	    indexes?: Index[];
	
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

