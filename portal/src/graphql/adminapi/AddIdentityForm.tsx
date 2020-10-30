import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useNavigate } from "react-router-dom";
import { ITextFieldProps } from "@fluentui/react";
import { Context } from "@oursky/react-messageformat";

import { UserQuery_node_User } from "./query/__generated__/UserQuery";
import { CreateIdentityFunction } from "./mutations/createIdentityMutation";
import NavigationBlockerDialog from "../../NavigationBlockerDialog";
import PasswordField, { localValidatePassword } from "../../PasswordField";
import ButtonWithLoading from "../../ButtonWithLoading";
import { PortalAPIAppConfig } from "../../types";
import { AuthenticatorType } from "./__generated__/globalTypes";
import { nonNullable } from "../../util/types";

import styles from "./AddIdentityForm.module.scss";

interface AddIdentityFormProps {
  className?: string;
  appConfig: PortalAPIAppConfig | null;
  user: UserQuery_node_User | null;
  password: string;
  onPasswordChange: ITextFieldProps["onChange"];
  passwordFieldErrorMessage?: string;
  loginIdKey: "username" | "email" | "phone";
  loginId: string;
  loginIdField: React.ReactNode;
  isFormModified: boolean;
  createIdentity: CreateIdentityFunction;
  creatingIdentity: boolean;
  onLocalErrorMessageChange: (errorMessage?: string) => void;
}

function determineIsPasswordRequired(
  user: UserQuery_node_User | null
): boolean {
  const authenticators =
    user?.authenticators?.edges
      ?.map((edge) => edge?.node?.type)
      .filter(nonNullable) ?? [];
  const hasPasswordAuthenticator = authenticators.includes(
    AuthenticatorType.PASSWORD
  );
  return !hasPasswordAuthenticator;
}

const AddIdentityForm: React.FC<AddIdentityFormProps> = function AddIdentityForm(
  props: AddIdentityFormProps
) {
  const {
    className,
    appConfig,
    user,
    password,
    onPasswordChange,
    passwordFieldErrorMessage,
    loginId,
    loginIdKey,
    loginIdField,
    isFormModified,
    createIdentity,
    creatingIdentity,
    onLocalErrorMessageChange,
  } = props;
  const navigate = useNavigate();
  const { renderToString } = useContext(Context);

  const isPasswordRequired = useMemo(() => {
    return determineIsPasswordRequired(user);
  }, [user]);

  const passwordPolicy = useMemo(() => {
    return appConfig?.authenticator?.password?.policy ?? {};
  }, [appConfig]);

  const [submittedForm, setSubmittedForm] = useState<boolean>(false);

  const onFormSubmit = useCallback(
    (ev: React.SyntheticEvent<HTMLElement>) => {
      ev.preventDefault();
      ev.stopPropagation();

      if (isPasswordRequired) {
        const localErrorMessageMap = localValidatePassword(
          renderToString,
          passwordPolicy,
          password
        );
        onLocalErrorMessageChange(localErrorMessageMap?.password);

        if (localErrorMessageMap != null) {
          return;
        }
      }

      createIdentity({ key: loginIdKey, value: loginId })
        .then((identity) => {
          if (identity != null) {
            setSubmittedForm(true);
          }
        })
        .catch(() => {});
    },
    [
      loginIdKey,
      loginId,
      createIdentity,
      isPasswordRequired,
      password,
      passwordPolicy,
      renderToString,
      onLocalErrorMessageChange,
    ]
  );

  useEffect(() => {
    if (submittedForm) {
      navigate("..#connected-identities");
    }
  }, [submittedForm, navigate]);

  return (
    <form className={className} onSubmit={onFormSubmit}>
      <NavigationBlockerDialog
        blockNavigation={!submittedForm && isFormModified}
      />
      {loginIdField}
      {isPasswordRequired && (
        <PasswordField
          className={styles.password}
          textFieldClassName={styles.passwordField}
          passwordPolicy={passwordPolicy}
          label={renderToString("AddUsernameScreen.password.label")}
          value={password}
          onChange={onPasswordChange}
          errorMessage={passwordFieldErrorMessage}
        />
      )}
      <ButtonWithLoading
        type="submit"
        disabled={!isFormModified || submittedForm}
        labelId="add"
        loading={creatingIdentity}
      />
    </form>
  );
};

export default AddIdentityForm;
