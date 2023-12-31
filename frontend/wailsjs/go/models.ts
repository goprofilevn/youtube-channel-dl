export namespace main {
	
	export class Folder {
	    folder: string;
	    err: any;
	
	    static createFrom(source: any = {}) {
	        return new Folder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.folder = source["folder"];
	        this.err = source["err"];
	    }
	}

}

