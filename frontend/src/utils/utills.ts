import { Ref } from "vue";

export function mapLabel(category: string, data: any): string {
  const values = Object.values(data) as any;
  switch (category) {
    case "database":
      return `${values[0].host}@${values[0].database}`;
    case "ssh":
      return "bi-terminal-fill";
    case "apiKeys":
      return `${Object.keys(values[0])[0]}`;
    default:
      return "bi-gear";
  }
}



function validateFiles(value:number) {
  
}
