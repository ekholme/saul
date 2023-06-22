library(tidyverse)

set.seed(0408)

fls <- dir("./data/sy22_spring", full.names = TRUE)

df_raw <- map_dfr(
    fls,
    read_csv
) |>
    janitor::clean_names()

schs <- unique(df_raw$sch_name)

schs_anon <- c("School A", "School B", "School C")

names(schs_anon) <- schs

# some filtering of tests and renaming schools

df_filtered <- df_raw |>
    filter(str_detect(test, "CAT")) |>
    mutate(
        sch_name = str_replace_all(sch_name, schs_anon),
        response = if_else(response == "COR", TRUE, FALSE)
    )

# summarize sdbq data
df_anon <- df_filtered |>
    group_by(sch_name, test, sol, item_descriptor) |>
    summarize(
        n_ask = n(),
        n_correct = sum(response, na.rm = TRUE),
        pct_correct = n_correct / n_ask,
        q = n_ask - n_correct
    ) |>
    ungroup()

df_sliced <- df_anon |>
    group_by(sch_name, test) |>
    slice_max(order_by = q, n = 3, with_ties = FALSE) |>
    select(sch_name, test, sol, item_descriptor, q) |>
    ungroup()

# getting a set of unique sols/item descriptors to get best practices for
sols_unique <- df_sliced |>
    distinct(test, sol, item_descriptor)

write_csv(sols_unique, "./data/standards_items_bp.csv")

#read in lookup
math_lookup <- read_csv("./data/math_items_bp_lookup.csv")

math_df <- df_sliced |>
    filter(str_detect(test, "Math")) |>
    select(-item_descriptor) |>
    left_join(math_lookup) |>
    select(sch_name, test, item_descriptor, best_practice, q)

write_csv(math_df, "./data/toy_math_data_bp.csv")
