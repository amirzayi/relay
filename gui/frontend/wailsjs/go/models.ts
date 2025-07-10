export namespace config {
	
	export class File {
	    Name: string;
	    Size: number;
	    Path: string;
	    Parents: string[];
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Size = source["Size"];
	        this.Path = source["Path"];
	        this.Parents = source["Parents"];
	    }
	}

}

