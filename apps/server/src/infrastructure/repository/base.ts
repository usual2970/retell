import { PrismaClient } from "@prisma/client";

// 定义 Prisma 相关的类型
type PrismaWhereCondition = Record<string, any>;
type PrismaOrderByCondition =
  | Record<string, "asc" | "desc" | Record<string, any>>
  | Array<Record<string, "asc" | "desc" | Record<string, any>>>;

// 定义操作符类型
type WhereOperator =
  | ""
  | "="
  | "eq"
  | "!="
  | "<>"
  | "ne"
  | "not"
  | ">"
  | "gt"
  | ">="
  | "gte"
  | "<"
  | "lt"
  | "<="
  | "lte"
  | "like"
  | "contains"
  | "not like"
  | "not contains"
  | "starts with"
  | "startswith"
  | "ends with"
  | "endswith"
  | "in"
  | "not in"
  | "notin"
  | "is null"
  | "isnull"
  | "is not null"
  | "isnotnull"
  | "between"
  | "not between";

// 定义排序方向类型
type SortDirection = "asc" | "desc";

// 定义排序对象类型
type OrderByObject = { [key: string]: SortDirection };

// 定义嵌套排序对象类型
type NestedOrderByObject = {
  [key: string]: SortDirection | NestedOrderByObject;
};

let prisma: PrismaClient | null = null;

export const getPrismaClient = (): PrismaClient => {
  if (!prisma) {
    prisma = new PrismaClient({
      log:
        process.env.NODE_ENV === "development"
          ? ["query", "info", "warn", "error"]
          : ["error"],
      errorFormat: "pretty",
    });

    // Handle connection errors
    prisma.$connect().catch((error) => {
      console.error("Failed to connect to database:", error);
      throw error;
    });
  }
  return prisma;
};

export const disconnectPrisma = async (): Promise<void> => {
  if (prisma) {
    await prisma.$disconnect();
    prisma = null;
  }
};

/**
 * 将 Map<string, any> 转换为 Prisma where 条件
 * 支持多种操作符和条件类型
 *
 * @param whereMap - 包含查询条件的 Map
 * @returns Prisma where 对象
 */
export const mapToWhereCondition = (
  whereMap: Map<string, unknown>
): PrismaWhereCondition => {
  const where: PrismaWhereCondition = {};

  for (const [key, value] of whereMap.entries()) {
    const parts = key.trim().split(/\s+/);
    const field = parts[0];
    const operator = parts.slice(1).join(" ").toLowerCase() as WhereOperator;

    if (!field) continue;

    switch (operator) {
      // 等于 (默认)
      case "":
      case "=":
      case "eq":
        where[field] = value;
        break;

      // 不等于
      case "!=":
      case "<>":
      case "ne":
      case "not":
        where[field] = { not: value };
        break;

      // 大于
      case ">":
      case "gt":
        where[field] = { gt: value };
        break;

      // 大于等于
      case ">=":
      case "gte":
        where[field] = { gte: value };
        break;

      // 小于
      case "<":
      case "lt":
        where[field] = { lt: value };
        break;

      // 小于等于
      case "<=":
      case "lte":
        where[field] = { lte: value };
        break;

      // 包含 (字符串)
      case "like":
      case "contains":
        if (typeof value === "string") {
          // 处理 SQL LIKE 模式
          if (value.startsWith("%") && value.endsWith("%")) {
            where[field] = { contains: value.slice(1, -1) };
          } else if (value.startsWith("%")) {
            where[field] = { endsWith: value.slice(1) };
          } else if (value.endsWith("%")) {
            where[field] = { startsWith: value.slice(0, -1) };
          } else {
            where[field] = { contains: value };
          }
        } else {
          where[field] = { contains: value };
        }
        break;

      // 不包含
      case "not like":
      case "not contains":
        if (typeof value === "string") {
          if (value.startsWith("%") && value.endsWith("%")) {
            where[field] = { not: { contains: value.slice(1, -1) } };
          } else if (value.startsWith("%")) {
            where[field] = { not: { endsWith: value.slice(1) } };
          } else if (value.endsWith("%")) {
            where[field] = { not: { startsWith: value.slice(0, -1) } };
          } else {
            where[field] = { not: { contains: value } };
          }
        } else {
          where[field] = { not: { contains: value } };
        }
        break;

      // 以...开始
      case "starts with":
      case "startswith":
        where[field] = { startsWith: value };
        break;

      // 以...结束
      case "ends with":
      case "endswith":
        where[field] = { endsWith: value };
        break;

      // 在数组中
      case "in":
        where[field] = { in: Array.isArray(value) ? value : [value] };
        break;

      // 不在数组中
      case "not in":
      case "notin":
        where[field] = { notIn: Array.isArray(value) ? value : [value] };
        break;

      // 为空
      case "is null":
      case "isnull":
        where[field] = null;
        break;

      // 不为空
      case "is not null":
      case "isnotnull":
        where[field] = { not: null };
        break;

      // 范围查询
      case "between":
        if (Array.isArray(value) && value.length === 2) {
          where[field] = { gte: value[0], lte: value[1] };
        }
        break;

      // 不在范围内
      case "not between":
        if (Array.isArray(value) && value.length === 2) {
          where[field] = { OR: [{ lt: value[0] }, { gt: value[1] }] };
        }
        break;

      // 默认处理 (当作等于)
      default:
        where[field] = value;
        break;
    }
  }

  return where;
};

/**
 * 构建复杂的 where 条件，支持 AND/OR 逻辑
 *
 * @param conditions - 条件数组
 * @param logic - 逻辑操作符 ('AND' | 'OR')
 * @returns Prisma where 对象
 */
export const buildComplexWhere = (
  conditions: Array<Map<string, unknown>>,
  logic: "AND" | "OR" = "AND"
): PrismaWhereCondition => {
  if (conditions.length === 0) return {};
  if (conditions.length === 1 && conditions[0])
    return mapToWhereCondition(conditions[0]);

  const whereConditions = conditions.map((condition) =>
    mapToWhereCondition(condition)
  );

  return {
    [logic]: whereConditions,
  };
};

/**
 * 从对象创建 where 条件的快捷方法
 *
 * @param obj - 包含查询条件的对象
 * @returns Prisma where 对象
 */
export const objectToWhere = (
  obj: Record<string, unknown>
): PrismaWhereCondition => {
  const whereMap = new Map(Object.entries(obj));
  return mapToWhereCondition(whereMap);
};

/**
 * 将排序数组转换为 Prisma orderBy 条件
 * 支持多字段排序和嵌套字段排序
 *
 * @param orderByArray - 排序条件数组
 * @returns Prisma orderBy 对象或数组
 */
export const arrayToOrderBy = (
  orderByArray: OrderByObject[]
): PrismaOrderByCondition | undefined => {
  if (!orderByArray || orderByArray.length === 0) {
    return undefined;
  }

  const processedOrderBy: NestedOrderByObject[] = orderByArray.map(
    (orderItem) => {
      const processedItem: NestedOrderByObject = {};

      for (const [key, direction] of Object.entries(orderItem)) {
        // 处理嵌套字段 (如 'user.name', 'category.priority')
        if (key.includes(".")) {
          const parts = key.split(".");
          let current: NestedOrderByObject = processedItem;

          // 构建嵌套对象
          for (let i = 0; i < parts.length - 1; i++) {
            const part = parts[i];
            if (part && !current[part]) {
              current[part] = {};
            }
            if (part) {
              current = current[part] as NestedOrderByObject;
            }
          }

          // 设置最终的排序方向
          const lastPart = parts[parts.length - 1];
          if (lastPart) {
            current[lastPart] = direction;
          }
        } else {
          // 处理普通字段
          processedItem[key] = direction;
        }
      }

      return processedItem;
    }
  );

  // 如果只有一个排序条件，返回对象；否则返回数组
  return processedOrderBy.length === 1 ? processedOrderBy[0] : processedOrderBy;
};

/**
 * 从字符串创建 orderBy 条件的快捷方法
 *
 * @param orderByString - 排序字符串，格式如 "createdAt desc, title asc"
 * @returns Prisma orderBy 对象或数组
 */
export const stringToOrderBy = (
  orderByString: string
): PrismaOrderByCondition | undefined => {
  if (!orderByString || orderByString.trim() === "") {
    return undefined;
  }

  const orderByArray: OrderByObject[] = orderByString
    .split(",")
    .map((item) => item.trim())
    .filter((item) => item !== "")
    .map((item) => {
      const parts = item.split(/\s+/);
      const field = parts[0];
      const direction: SortDirection =
        parts[1]?.toLowerCase() === "desc" ? "desc" : "asc";

      return field ? { [field]: direction } : {};
    })
    .filter((item) => Object.keys(item).length > 0);

  return arrayToOrderBy(orderByArray);
};

/**
 * 验证 orderBy 字段是否合法
 *
 * @param orderByArray - 排序条件数组
 * @param allowedFields - 允许的字段列表
 * @returns 验证结果
 */
export const validateOrderBy = (
  orderByArray: OrderByObject[],
  allowedFields: string[]
): { valid: boolean; invalidFields: string[] } => {
  const invalidFields: string[] = [];

  for (const orderItem of orderByArray) {
    for (const field of Object.keys(orderItem)) {
      // 处理嵌套字段
      const baseField = field.split(".")[0];
      if (
        baseField &&
        !allowedFields.includes(baseField) &&
        !allowedFields.includes(field)
      ) {
        invalidFields.push(field);
      }
    }
  }

  return {
    valid: invalidFields.length === 0,
    invalidFields,
  };
};

/**
 * 合并多个 orderBy 条件
 *
 * @param orderByArrays - 多个排序条件数组
 * @returns 合并后的 Prisma orderBy 条件
 */
export const mergeOrderBy = (
  ...orderByArrays: OrderByObject[][]
): PrismaOrderByCondition | undefined => {
  const mergedArray = orderByArrays
    .flat()
    .filter((item) => item && Object.keys(item).length > 0);
  return arrayToOrderBy(mergedArray);
};
