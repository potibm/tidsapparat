import CategoryIcon from "@mui/icons-material/Category";
import CategoriesList from "./CategoriesList";
import CategoriesCreate from "./CategoriesCreate";
import CategoriesEdit from "./CategoriesEdit";

export default {
  name: "categories",
  options: { label: "Categories" },
  list: CategoriesList,
  create: CategoriesCreate,
  edit: CategoriesEdit,
  icon: CategoryIcon,
};
