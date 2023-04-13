library(tidyverse)

df_raw <- readRDS(here::here("data/winter_22_working.Rds"))

schs <- unique(df_raw$sch_name[df_raw$tst_yr == 2022 & df_raw$tst_window == "Winter"])

sch_sub <- sample(schs, size = 3)

df_sub <- subset(df_raw, sch_name %in% sch_sub)

# anonymize school names
schs_anon <- c("School A", "School B", "School C")

names(schs_anon) <- sch_sub

# replace with anon names
# filter to toy subset
# and calculate q
df_anon <- df_sub |>
    mutate(
        sch_name = str_replace_all(sch_name, schs_anon),
        q = (1 - pct_correct) * n_ask
        ) |>
    filter(
        tst_window == "Winter",
        tst_yr == 2022,
        st_group == "All"
    ) |>
    select(sch_name, test, item_descriptor, q)

#write out 
write_csv(df_anon, here::here("data/toy_data.csv"), col_names = FALSE)

# write out a df with just the unique tests
tsts_df <- df_anon |>
    distinct(sch_name, test)

write_csv(tsts_df, here::here("data/tsts.csv"), col_names = FALSE)